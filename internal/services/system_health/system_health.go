package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/iRankHub/backend/internal/utils"
)

type SystemHealthService struct {
	clientset     *kubernetes.Clientset
	metricsClient *versioned.Clientset
}

type SystemMetrics struct {
	CPUUsagePercentage         float64
	MemoryUsagePercentage      float64
	EphemeralStoragePercentage float64
	PVCStoragePercentage       float64
	NodeCount                  int
	PodCount                   int
	PVCCount                   int
}

func NewSystemHealthService() (*SystemHealthService, error) {
	var config *rest.Config
	var err error

	// Check for in-cluster configuration first
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get in-cluster config: %v", err)
		}
	} else {
		// Use out-of-cluster configuration
		host := os.Getenv("KUBERNETES_SERVICE_HOST")
		port := os.Getenv("KUBERNETES_SERVICE_PORT")

		if host != "" && port != "" {
			config = &rest.Config{
				Host: fmt.Sprintf("https://%s:%s", host, port),
			}
		} else {
			kubeconfigPath := os.Getenv("KUBECONFIG")
			if kubeconfigPath == "" {
				kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
			}
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
			if err != nil {
				return nil, fmt.Errorf("failed to get Kubernetes config: %v", err)
			}
		}
	}

	// Skip TLS verification
	config.Insecure = true
	config.TLSClientConfig = rest.TLSClientConfig{Insecure: true}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes clientset: %v", err)
	}

	metricsClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Metrics clientset: %v", err)
	}

	return &SystemHealthService{
		clientset:     clientset,
		metricsClient: metricsClient,
	}, nil
}

func (s *SystemHealthService) GetSystemHealth(ctx context.Context, token string) (*SystemMetrics, error) {
	if err := s.validateAdminRole(token); err != nil {
		return nil, err
	}

	metrics := &SystemMetrics{}

	// Fetch node metrics
	nodeMetrics, err := s.metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch node metrics: %v", err)
	}

	var totalCPUUsage, totalMemoryUsage int64
	for _, nodeMetric := range nodeMetrics.Items {
		totalCPUUsage += nodeMetric.Usage.Cpu().MilliValue()
		totalMemoryUsage += nodeMetric.Usage.Memory().Value()
	}

	// Fetch nodes
	nodes, err := s.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nodes: %v", err)
	}

	metrics.NodeCount = len(nodes.Items)

	var totalCPUCapacity, totalMemoryCapacity, totalEphemeralStorageCapacity int64
	for _, node := range nodes.Items {
		totalCPUCapacity += node.Status.Capacity.Cpu().MilliValue()
		totalMemoryCapacity += node.Status.Capacity.Memory().Value()
		totalEphemeralStorageCapacity += node.Status.Capacity.StorageEphemeral().Value()
	}

	metrics.CPUUsagePercentage = float64(totalCPUUsage) / float64(totalCPUCapacity) * 100
	metrics.MemoryUsagePercentage = float64(totalMemoryUsage) / float64(totalMemoryCapacity) * 100

	// Fetch pods
	pods, err := s.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pods: %v", err)
	}

	metrics.PodCount = len(pods.Items)

	var totalEphemeralStorageUsage int64
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if container.Resources.Requests != nil {
				totalEphemeralStorageUsage += container.Resources.Requests.StorageEphemeral().Value()
			}
		}
	}

	metrics.EphemeralStoragePercentage = float64(totalEphemeralStorageUsage) / float64(totalEphemeralStorageCapacity) * 100

	// Fetch PVCs
	pvcs, err := s.clientset.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PVCs: %v", err)
	}

	metrics.PVCCount = len(pvcs.Items)

	var totalPVCCapacity, totalPVCUsage int64
	for _, pvc := range pvcs.Items {
		totalPVCCapacity += pvc.Spec.Resources.Requests.Storage().Value()
		if pvc.Status.Phase == corev1.ClaimBound {
			totalPVCUsage += pvc.Status.Capacity.Storage().Value()
		}
	}

	if totalPVCCapacity > 0 {
		metrics.PVCStoragePercentage = float64(totalPVCUsage) / float64(totalPVCCapacity) * 100
	}

	return metrics, nil
}

func (s *SystemHealthService) validateAdminRole(token string) error {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return fmt.Errorf("unauthorized: only admins can perform this action")
	}

	return nil
}
