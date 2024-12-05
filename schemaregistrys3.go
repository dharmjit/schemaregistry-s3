package schemaregistrys3

import (
	"context"
	"fmt"

	"github.com/dharmjit/schemaregistry-s3/store"
	"github.com/dharmjit/schemaregistry-s3/validator"
)

type schemeFormat string

const (
	AvroSchemeFormat     schemeFormat = "avro"
	JsonSchemeFormat     schemeFormat = "json"
	ProtobufSchemeFormat schemeFormat = "protobuf"
)

type S3SchemaRegistry struct {
	store     store.SchemaStore
	validator validator.Validator
}

func NewS3SchemaRegistry(ctx context.Context, s3storageprovider string, region, bucket string, schemaFormat schemeFormat) (*S3SchemaRegistry, error) {
	var storeProvider store.SchemaStore
	switch s3storageprovider {
	case "minio":
		miniostore, err := store.NewMinioStore(ctx, "127.0.0.1:8000", region, bucket)
		if err != nil {
			return nil, err
		}
		storeProvider = miniostore
	case "s3":
		s3store, err := store.NewS3Store(ctx, region, bucket)
		if err != nil {
			return nil, err
		}
		storeProvider = s3store
	}

	var schemaValidator validator.Validator
	switch schemaFormat {
	case JsonSchemeFormat:
		schemaValidator = validator.NewJsonSchemaSpecValidator()
	}
	return &S3SchemaRegistry{
		store:     storeProvider,
		validator: schemaValidator,
	}, nil
}

func (s *S3SchemaRegistry) PutSchema(ctx context.Context, schema, version string, schemaData []byte) error {
	err := s.validator.ValidateSchema(schemaData)
	if err != nil {
		return fmt.Errorf("failed to validate schema: %w", err)
	}
	return s.store.PutSchema(ctx, schema, version, schemaData)
}

func (s *S3SchemaRegistry) ListSchemas(ctx context.Context) ([]string, error) {
	return s.store.ListSchemas(ctx)
}

func (s *S3SchemaRegistry) ListSchemaVersions(ctx context.Context, schema string) ([]string, error) {
	return s.store.ListSchemaVersions(ctx, schema)
}

func (s *S3SchemaRegistry) GetSchema(ctx context.Context, schema, version string) ([]byte, error) {
	return s.store.GetSchema(ctx, schema, version)
}

func (s *S3SchemaRegistry) DeleteSchema(ctx context.Context, schema, version string) error {
	return s.store.DeleteSchema(ctx, schema, version)
}
