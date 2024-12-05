package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
)

func setupMinioTestClient() *minio.Client {
	minioClient, err := minio.New("127.0.0.1:8000", &minio.Options{
		Creds:  credentials.NewStaticV2("minioadmin", "minioadmin", ""),
		Region: "us-east-1",
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
	return minioClient
}

func setupMinioBucket(client *minio.Client, bucketName string) error {
	err := client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	// if the bucket already exists, ignore the error
	if err != nil {
		if minio.ToErrorResponse(err).Code == "BucketAlreadyOwnedByYou" || minio.ToErrorResponse(err).Code == "BucketAlreadyExists" {
			fmt.Printf("Bucket %s already exists, ignoring error\n", bucketName)
			return nil
		}
		return err
	}
	return nil
}

func deleteMinioBucket(ctx context.Context, client *minio.Client, bucketName string) error {
	for object := range client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix: "schemas/",
	}) {
		client.RemoveObject(ctx, bucketName, object.Key, minio.RemoveObjectOptions{})
	}
	return nil
}

func TestMinioPutSchema(t *testing.T) {
	ctx := context.Background()
	minioClient := setupMinioTestClient()
	store := &MinioStore{
		Client:    minioClient,
		Bucket:    "minio-test-bucket",
		Prefix:    "schemas/",
		Delimiter: "/",
	}
	err := deleteMinioBucket(ctx, minioClient, store.Bucket)
	assert.NoError(t, err)
	err = setupMinioBucket(minioClient, store.Bucket)
	assert.NoError(t, err)

	schema1 := "test-schema1"
	schema2 := "test-schema2"
	versionv1 := "v1"
	versionv2 := "v2"
	schemaData := []byte(`{"type": "record", "name": "TestSchema", "fields": [{"name": "field1", "type": "string"}]}`)

	err = store.PutSchema(ctx, schema1, versionv1, schemaData)
	assert.NoError(t, err)
	err = store.PutSchema(ctx, schema1, versionv2, schemaData)
	assert.NoError(t, err)

	err = store.PutSchema(ctx, schema2, versionv1, schemaData)
	assert.NoError(t, err)
	err = store.PutSchema(ctx, schema2, versionv2, schemaData)
	assert.NoError(t, err)

	schemas, err := store.ListSchemas(ctx)
	assert.NoError(t, err)
	assert.Len(t, schemas, 2)
	assert.Contains(t, schemas, schema1)
	assert.Contains(t, schemas, schema2)
	versions, err := store.ListSchemaVersions(ctx, schema1)
	assert.Len(t, versions, 2)
	assert.NoError(t, err)
	assert.Contains(t, versions, versionv1)
	assert.Contains(t, versions, versionv2)

	versions, err = store.ListSchemaVersions(ctx, schema2)
	assert.Len(t, versions, 2)
	assert.NoError(t, err)
	assert.Contains(t, versions, versionv1)
	assert.Contains(t, versions, versionv2)
}
