package schemaregistrys3

import (
	"context"
	"os"
	"testing"
)

func TestIntegrationS3SchemaRegistry_PutSchema(t *testing.T) {
	bucketName := "test-integration-bucket"
	os.Setenv("AWS_ACCESS_KEY", "minioadmin")
	os.Setenv("AWS_SECRET_KEY", "minioadmin")
	registry, err := NewS3SchemaRegistry(context.Background(), "minio", "us-east-1", bucketName, JsonSchemeFormat)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	err = registry.PutSchema(ctx,
		"eventtype1",
		"v1",
		[]byte(`{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"type": "object",
		"properties": {
			"name": { "type": "string" },
			"age": { "type": "integer", "minimum": 0 }
		},
		"required": ["name", "age"]
	}`))
	if err != nil {
		t.Fatal(err)
	}
}
