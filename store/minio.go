package store

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStore struct {
	// Minio client
	Client *minio.Client
	// Bucket name
	Bucket string
	// Prefix name
	Prefix string
	// Delimiter Character
	Delimiter string
}

// NewS3Store creates a new S3Store.
func NewMinioStore(ctx context.Context, endpoint, region, bucket string) (*MinioStore, error) {
	client, err := setupMinioClient(ctx, endpoint, region)
	if err != nil {
		return nil, err
	}
	return &MinioStore{
		Client:    client,
		Bucket:    bucket,
		Prefix:    "schemas/",
		Delimiter: "/",
	}, nil
}

func setupMinioClient(_ context.Context, endpoint, region string) (*minio.Client, error) {
	// Create a minio client
	creds := credentials.NewEnvAWS()
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  creds,
		Region: region,
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

// PutSchema puts the schema for a given schema and version.
func (s *MinioStore) PutSchema(ctx context.Context, schema, version string, schemaData []byte) error {
	_, err := s.Client.PutObject(ctx, s.Bucket, s.getSchemaVersionKey(schema, version), bytes.NewReader(schemaData), int64(len(schemaData)), minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}
	return nil
}

// ListSchemas returns the list of schemas in the store.
func (s *MinioStore) ListSchemas(ctx context.Context) ([]string, error) {
	uniqueSchemas := make(map[string]bool)
	for object := range s.Client.ListObjects(ctx, s.Bucket, minio.ListObjectsOptions{Prefix: s.Prefix}) {
		// Extract the schema name part of the key
		// Assuming the format: "schemas/<schema_name>/..."
		parts := strings.Split(object.Key, s.Delimiter)
		if len(parts) > 1 {
			schemaName := parts[1]
			uniqueSchemas[schemaName] = true
		}
	}
	schemas := make([]string, 0, len(uniqueSchemas))
	for schema := range uniqueSchemas {
		schemas = append(schemas, schema)
	}
	return schemas, nil
}

// ListSchemaVersions returns the list of versions for a given schema.
func (s *MinioStore) ListSchemaVersions(ctx context.Context, schema string) ([]string, error) {
	var objects []string
	for object := range s.Client.ListObjects(ctx, s.Bucket, minio.ListObjectsOptions{
		Prefix: s.Prefix + schema + s.Delimiter,
	}) {
		objects = append(objects, strings.Split(object.Key, s.Delimiter)[2])
	}
	return objects, nil
}

// GetSchema returns the schema for a given schema and version.
func (s *MinioStore) GetSchema(ctx context.Context, schema, version string) ([]byte, error) {
	return nil, nil
}

// DeleteSchema deletes the schema for a given schema and version.
func (s *MinioStore) DeleteSchema(ctx context.Context, schema, version string) error {
	return nil
}

func (s *MinioStore) getSchemaVersionKey(schema, version string) string {
	return *aws.String(s.Prefix + schema + s.Delimiter + version + s.Delimiter + schema + "_" + version + ".json")
}
