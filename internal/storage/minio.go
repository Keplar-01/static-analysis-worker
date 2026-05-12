package storage

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	bucketSourceCodes      = "source-codes"
	bucketAnalysisArtifact = "analysis-artifacts"
)

type MinIOClient struct {
	client *minio.Client
}

func NewMinIOClient(endpoint, accessKey, secretKey string) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("minio connect: %w", err)
	}
	return &MinIOClient{client: client}, nil
}

func (m *MinIOClient) DownloadSource(ctx context.Context, s3Path, workDir string) (string, error) {
	bucket, object := splitS3Path(s3Path, bucketSourceCodes)

	dst := filepath.Join(workDir, filepath.Base(object))
	if err := m.client.FGetObject(ctx, bucket, object, dst, minio.GetObjectOptions{}); err != nil {
		return "", fmt.Errorf("download %s/%s: %w", bucket, object, err)
	}
	return dst, nil
}

func (m *MinIOClient) UploadArtifact(ctx context.Context, taskID string, payload []byte) (string, error) {
	if err := m.ensureBucket(ctx, bucketAnalysisArtifact); err != nil {
		return "", err
	}

	objectKey := fmt.Sprintf("%s/static-out.json", taskID)

	_, err := m.client.PutObject(
		ctx,
		bucketAnalysisArtifact,
		objectKey,
		bytes.NewReader(payload),
		int64(len(payload)),
		minio.PutObjectOptions{ContentType: "application/json"},
	)
	if err != nil {
		return "", fmt.Errorf("put %s/%s: %w", bucketAnalysisArtifact, objectKey, err)
	}

	return fmt.Sprintf("%s/%s", bucketAnalysisArtifact, objectKey), nil
}

func (m *MinIOClient) ensureBucket(ctx context.Context, bucket string) error {
	ok, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("bucket exists %s: %w", bucket, err)
	}
	if ok {
		return nil
	}
	if err := m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("make bucket %s: %w", bucket, err)
	}
	return nil
}

func splitS3Path(s3Path, defaultBucket string) (string, string) {
	parts := strings.SplitN(s3Path, "/", 2)
	if len(parts) == 2 && looksLikeBucket(parts[0]) {
		return parts[0], parts[1]
	}
	return defaultBucket, s3Path
}

func looksLikeBucket(s string) bool {
	switch s {
	case bucketSourceCodes, bucketAnalysisArtifact:
		return true
	}
	return false
}
