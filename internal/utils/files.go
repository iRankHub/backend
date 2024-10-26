package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client   *s3.Client
	bucket   string
	endpoint string
	region   string
}

func NewS3Client() (*S3Client, error) {
	region := os.Getenv("S3_REGION")
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	bucket := os.Getenv("S3_BUCKET")

	if endpoint == "" || region == "" || accessKey == "" || secretKey == "" || bucket == "" {
		return nil, fmt.Errorf("missing required S3 configuration")
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               fmt.Sprintf("https://%s", endpoint),
			HostnameImmutable: true,
			SigningRegion:     region,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load S3 config: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3Client{
		client:   client,
		bucket:   bucket,
		endpoint: endpoint,
		region:   region,
	}, nil
}

func (s *S3Client) UploadFile(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	// Ensure the key is URL safe
	safeKey := strings.ReplaceAll(key, " ", "-")
	safeKey = strings.ReplaceAll(safeKey, "%", "-")
	safeKey = url.PathEscape(safeKey)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(safeKey),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	// Construct URL in the same format as the frontend
	fileURL := fmt.Sprintf("https://%s.%s.linodeobjects.com/%s",
		s.bucket, s.region, safeKey)

	return fileURL, nil
}

func (s *S3Client) GetSignedURL(ctx context.Context, fileURL string, expires time.Duration) (string, error) {
	// Extract the key from the URL
	key := ExtractKeyFromURL(fileURL)
	if key == "" {
		return "", fmt.Errorf("invalid file URL")
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	presignClient := s3.NewPresignClient(s.client)
	presignedURL, err := presignClient.PresignGetObject(ctx, input, s3.WithPresignExpires(expires))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return presignedURL.URL, nil
}

func (s *S3Client) DownloadFile(ctx context.Context, fileURL string) ([]byte, error) {
	key := ExtractKeyFromURL(fileURL)
	if key == "" {
		return nil, fmt.Errorf("invalid file URL")
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	result, err := s.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func (s *S3Client) DeleteFile(ctx context.Context, fileURL string) error {
	key := ExtractKeyFromURL(fileURL)
	if key == "" {
		return fmt.Errorf("invalid file URL")
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

// ExtractKeyFromURL extracts the object key from various URL formats
func ExtractKeyFromURL(fileURL string) string {
	if fileURL == "" {
		return ""
	}

	// Parse the URL
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return ""
	}

	// Clean the path
	cleanPath := path.Clean(parsedURL.Path)

	// Split the path into segments
	segments := strings.Split(cleanPath, "/")

	// Remove empty segments and the bucket name if present
	var validSegments []string
	for _, segment := range segments {
		if segment != "" && segment != "irankhub-bucket" {
			validSegments = append(validSegments, segment)
		}
	}

	// Join the remaining segments
	if len(validSegments) > 0 {
		return strings.Join(validSegments, "/")
	}

	return ""
}
