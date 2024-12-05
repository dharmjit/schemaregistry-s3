package store

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
)

func setupS3TestClient() *s3.Client {
	return s3.NewFromConfig(aws.Config{Region: "us-east-1"}, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://127.0.0.1:8000")
		o.Credentials = credentials.NewStaticCredentialsProvider("minioadmin", "minioadmin", "")
	})
}

func setupS3Bucket(client *s3.Client, bucketName string) error {
	_, err := client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	// if the bucket already exists, ignore the error
	var bae *types.BucketAlreadyExists
	var boe *types.BucketAlreadyOwnedByYou
	if err != nil && !errors.As(err, &bae) && !errors.As(err, &boe) {
		return err
	}
	return nil
}

func deleteS3Bucket(ctx context.Context, client *s3.Client, bucketName string) error {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String("schemas/"),
	}
	paginator := s3.NewListObjectsV2Paginator(client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list objects: %w", err)
		}
		if len(output.Contents) == 0 {
			continue
		}
		objectIds := make([]types.ObjectIdentifier, len(output.Contents))
		for i, object := range output.Contents {
			objectIds[i] = types.ObjectIdentifier{
				Key: object.Key,
			}
		}
		deleteInput := &s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &types.Delete{
				Objects: objectIds,
				Quiet:   aws.Bool(false),
			},
		}
		_, err = client.DeleteObjects(ctx, deleteInput)
		if err != nil {
			return fmt.Errorf("failed to delete objects: %w", err)
		}
	}
	return nil
}

func TestS3PutSchema(t *testing.T) {
	ctx := context.Background()
	s3Client := setupS3TestClient()
	store := &S3Store{
		Client:    s3Client,
		Bucket:    "s3-test-bucket",
		Prefix:    "schemas/",
		Delimiter: "/",
	}
	err := deleteS3Bucket(ctx, s3Client, store.Bucket)
	assert.NoError(t, err)
	err = setupS3Bucket(s3Client, store.Bucket)
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
