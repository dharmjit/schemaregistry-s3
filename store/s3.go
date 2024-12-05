package store

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Store struct {
	// Minio client
	Client *s3.Client
	// Bucket name
	Bucket string
	// Prefix name
	Prefix string
	// Delimiter Character
	Delimiter string
}

// NewS3Store creates a new S3Store.
func NewS3Store(ctx context.Context, region, bucket string) (*S3Store, error) {
	client, err := setupClient(ctx, region)
	if err != nil {
		return nil, err
	}
	return &S3Store{
		Client:    client,
		Bucket:    bucket,
		Prefix:    "schemas/",
		Delimiter: "/",
	}, nil
}

func setupClient(ctx context.Context, region string) (*s3.Client, error) {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg)
	return s3Client, nil
}

// PutSchema puts the schema for a given schema and version.
func (s *S3Store) PutSchema(ctx context.Context, schema, version string, schemaData []byte) error {
	// Create the PutObjectInput
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s.getSchemaVersionKey(schema, version)),
		Body:   bytes.NewReader(schemaData),
	}
	_, err := s.Client.PutObject(ctx, input)
	return err
}

// ListSchemas returns the list of schemas in the store.
func (s *S3Store) ListSchemas(ctx context.Context) ([]string, error) {
	output, err := s.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.Bucket),
		Prefix: aws.String(s.Prefix),
	})
	if err != nil {
		return nil, err
	}
	uniqueSchemas := make(map[string]bool)
	for _, obj := range output.Contents {
		// Extract the schema name part of the key
		// Assuming the format: "schemas/<schema_name>/..."
		parts := strings.Split(*obj.Key, s.Delimiter)
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
func (s *S3Store) ListSchemaVersions(ctx context.Context, schema string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.Bucket),
		Prefix: aws.String(s.Prefix + schema + s.Delimiter),
	}
	output, err := s.Client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, err
	}

	var versions []string
	for _, item := range output.Contents {
		versions = append(versions, strings.Split(*item.Key, s.Delimiter)[2])
	}
	return versions, nil
}

// GetSchema returns the schema for a given schema and version.
func (s *S3Store) GetSchema(ctx context.Context, schema, version string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s.getSchemaVersionKey(schema, version)),
	}
	result, err := s.Client.GetObject(ctx, input)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	// Read the object from the response
	data, err := io.ReadAll(result.Body)
	return data, err
}

// DeleteSchema deletes the schema for a given schema and version.
func (s *S3Store) DeleteSchema(ctx context.Context, schema, version string) error {
	input := &s3.DeleteObjectInput{
		Bucket: &s.Bucket,
		Key:    &schema,
	}
	_, err := s.Client.DeleteObject(ctx, input)
	return err
}

func (s *S3Store) getSchemaVersionKey(schema, version string) string {
	return *aws.String(s.Prefix + schema + s.Delimiter + version + s.Delimiter + schema + "_" + version + ".json")
}
