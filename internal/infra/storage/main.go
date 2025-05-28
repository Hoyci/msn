package storage

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageClient struct {
	client *minio.Client
}

func NewStorageClient(endpoint, accessKey, secretKey string) *StorageClient {
	storageClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln("failed to create minio client: ", err)
	}

	return &StorageClient{client: storageClient}
}

func (c *StorageClient) UploadFile(bucketName, objectName string, fileHeader *multipart.FileHeader) (string, error) {
	ctx := context.Background()
	exists, err := c.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", fmt.Errorf("error checking bucket %s existence: %w", bucketName, err)
	}
	if !exists {
		err = c.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	info, err := c.client.PutObject(ctx, bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object %w", err)
	}

	return info.Key, nil
}
