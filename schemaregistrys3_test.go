package schemaregistrys3

import (
	"context"
	"testing"

	mockschemastore "github.com/dharmjit/schemaregistry-s3/store/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestS3SchemaRegistry_ListSchemas(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mockschemastore.NewMockSchemaStore(ctrl)
	mockStore.EXPECT().ListSchemas(gomock.Any()).Return([]string{"schema1", "schema2"}, nil)

	registry := S3SchemaRegistry{store: mockStore}

	ctx := context.Background()
	schemas, err := registry.ListSchemas(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []string{"schema1", "schema2"}, schemas)
}

func TestS3SchemaRegistry_ListSchemaVersions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mockschemastore.NewMockSchemaStore(ctrl)
	mockStore.EXPECT().ListSchemaVersions(gomock.Any(), "schema1").Return([]string{"v1", "v2"}, nil)

	registry := S3SchemaRegistry{store: mockStore}

	ctx := context.Background()
	versions, err := registry.ListSchemaVersions(ctx, "schema1")
	assert.NoError(t, err)
	assert.Equal(t, []string{"v1", "v2"}, versions)
}

func TestS3SchemaRegistry_PutSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mockschemastore.NewMockSchemaStore(ctrl)
	mockStore.EXPECT().PutSchema(gomock.Any(), "schema1", "v1", []byte("schema data")).Return(nil)

	registry := S3SchemaRegistry{store: mockStore}

	ctx := context.Background()
	err := registry.PutSchema(ctx, "schema1", "v1", []byte("schema data"))
	assert.NoError(t, err)
}

func TestS3SchemaRegistry_GetSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mockschemastore.NewMockSchemaStore(ctrl)
	mockStore.EXPECT().GetSchema(gomock.Any(), "schema1", "v1").Return([]byte("schema data"), nil)

	registry := S3SchemaRegistry{store: mockStore}

	ctx := context.Background()
	data, err := registry.GetSchema(ctx, "schema1", "v1")
	assert.NoError(t, err)
	assert.Equal(t, []byte("schema data"), data)
}

func TestS3SchemaRegistry_DeleteSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mockschemastore.NewMockSchemaStore(ctrl)
	mockStore.EXPECT().DeleteSchema(gomock.Any(), "schema1", "v1").Return(nil)

	registry := S3SchemaRegistry{store: mockStore}

	ctx := context.Background()
	err := registry.DeleteSchema(ctx, "schema1", "v1")
	assert.NoError(t, err)
}
