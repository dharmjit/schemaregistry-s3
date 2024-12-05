package store

import (
	"context"
)

// SchemaStore is the interface that wraps the basic methods to interact with an S3 compatible store.
//
//go:generate go run go.uber.org/mock/mockgen -package mockschemastore -typed=true -source=store.go -destination ./mocks/store.go
type SchemaStore interface {
	// ListSchemas returns the list of schemas in the store.
	ListSchemas(ctx context.Context) ([]string, error)

	// ListSchemaVersions returns the list of versions for a given schema.
	ListSchemaVersions(ctx context.Context, schema string) ([]string, error)

	// PutSchema puts the schema for a given schema and version.
	PutSchema(ctx context.Context, schema, version string, schemaData []byte) error

	// GetSchema returns the schema for a given schema and version.
	GetSchema(ctx context.Context, schema, version string) ([]byte, error)

	// DeleteSchema deletes the schema for a given schema and version.
	DeleteSchema(ctx context.Context, schema, version string) error
}
