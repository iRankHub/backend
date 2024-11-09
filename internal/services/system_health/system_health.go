package services

import (
	"context"
	"encoding/json"
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

type StorageMetric struct {
	Timestamp  string `json:"timestamp"`
	UsedBytes  int64  `json:"used_bytes"`
	TotalBytes int64  `json:"total_bytes"`
}

type SystemMetrics struct {
	CPUUsagePercentage         float64
	MemoryUsagePercentage      float64
	EphemeralStorageUsed       int64
	EphemeralStorageTotal      int64
	EphemeralStoragePercentage float64
	PVCStorageUsed             int64
	PVCStorageTotal            int64
	PVCStoragePercentage       float64
	NodeCount                  int
	PodCount                   int
	PVCCount                   int
}

type NodeStats struct {
	Node struct {
		NodeName string `json:"nodeName"`
		Fs       struct {
			Time           string `json:"time"`
			AvailableBytes int64  `json:"availableBytes"`
			CapacityBytes  int64  `json:"capacityBytes"`
			UsedBytes      int64  `json:"usedBytes"`
		} `json:"fs"`
	} `json:"node"`
	Pods []struct {
		PodRef struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
		} `json:"podRef"`
		EphemeralStorage struct {
			Time           string `json:"time"`
			UsedBytes      int64  `json:"usedBytes"`
			AvailableBytes int64  `json:"availableBytes"`
			CapacityBytes  int64  `json:"capacityBytes"`
		} `json:"ephemeral-storage"`
	} `json:"pods"`
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

func (s *SystemHealthService) getPodMetrics(ctx context.Context, pod *corev1.Pod) (*StorageMetric, error) {
	// Execute command in the storage-metrics container to read the last line of metrics
	cmd := []string{"/bin/sh", "-c", "tail -n 1 /metrics/storage.log"}
	stdout, stderr, err := utils.ExecInPod(ctx, s.clientset, pod.Namespace, pod.Name, "storage-metrics", cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to exec in pod: %v, stderr: %s", err, stderr)
	}

	var metric StorageMetric
	if err := json.Unmarshal([]byte(stdout), &metric); err != nil {
		return nil, fmt.Errorf("failed to parse metrics json: %v", err)
	}

	return &metric, nil
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

	// Calculate total node capacities
	var totalCPUCapacity, totalMemoryCapacity int64
	for _, node := range nodes.Items {
		totalCPUCapacity += node.Status.Capacity.Cpu().MilliValue()
		totalMemoryCapacity += node.Status.Capacity.Memory().Value()
	}

	metrics.CPUUsagePercentage = float64(totalCPUUsage) / float64(totalCPUCapacity) * 100
	metrics.MemoryUsagePercentage = float64(totalMemoryUsage) / float64(totalMemoryCapacity) * 100

	// Get node stats for ephemeral storage
	metrics.EphemeralStorageUsed = 0
	metrics.EphemeralStorageTotal = 0

	for _, node := range nodes.Items {
		path := fmt.Sprintf("/api/v1/nodes/%s/proxy/stats/summary", node.Name)
		result, err := s.clientset.CoreV1().RESTClient().Get().AbsPath(path).DoRaw(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get node stats: %v", err)
		}

		var nodeStats NodeStats
		if err := json.Unmarshal(result, &nodeStats); err != nil {
			return nil, fmt.Errorf("failed to parse node stats: %v", err)
		}

		// Add node's total capacity
		metrics.EphemeralStorageTotal += nodeStats.Node.Fs.CapacityBytes

		// Sum up used storage across all pods
		for _, pod := range nodeStats.Pods {
			metrics.EphemeralStorageUsed += pod.EphemeralStorage.UsedBytes
		}
	}

	if metrics.EphemeralStorageTotal > 0 {
		metrics.EphemeralStoragePercentage = float64(metrics.EphemeralStorageUsed) / float64(metrics.EphemeralStorageTotal) * 100
	}

	// Fetch pods
	pods, err := s.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pods: %v", err)
	}

	metrics.PodCount = len(pods.Items)

	// Fetch PVCs and their actual usage from sidecar metrics
	pvcs, err := s.clientset.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PVCs: %v", err)
	}

	metrics.PVCCount = len(pvcs.Items)
	metrics.PVCStorageUsed = 0
	metrics.PVCStorageTotal = 0

	// Create a map of PVC name to pod for easy lookup
	pvcToPod := make(map[string]*corev1.Pod)
	for _, pod := range pods.Items {
		for _, volume := range pod.Spec.Volumes {
			if volume.PersistentVolumeClaim != nil {
				pvcToPod[volume.PersistentVolumeClaim.ClaimName] = pod.DeepCopy()
			}
		}
	}

	for _, pvc := range pvcs.Items {
		metrics.PVCStorageTotal += pvc.Spec.Resources.Requests.Storage().Value()

		if pod, exists := pvcToPod[pvc.Name]; exists {
			if storageMetric, err := s.getPodMetrics(ctx, pod); err == nil {
				metrics.PVCStorageUsed += storageMetric.UsedBytes
			}
		}
	}

	if metrics.PVCStorageTotal > 0 {
		metrics.PVCStoragePercentage = float64(metrics.PVCStorageUsed) / float64(metrics.PVCStorageTotal) * 100
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
