package validator

import (
	"testing"
)

func TestValidateSchema(t *testing.T) {
	validator := NewJsonSchemaSpecValidator()

	tests := []struct {
		name      string
		schemaDoc []byte
		wantErr   bool
	}{
		{
			name: "Valid schema",
			schemaDoc: []byte(`{
				"$schema": "http://json-schema.org/draft/2020-12/schema#",
				"type": "object",
				"properties": {
					"name": {
						"type": "string"
					}
				},
				"required": ["name"]
			}`),
			wantErr: false,
		},
		{
			name: "Invalid schema",
			schemaDoc: []byte(`{
				"$schema": "http://json-schema.org/draft/2020-12/schema#",
				"type": "object",
				"properties": {
					"name": {
						"type": "invalid-type"
					}
				},
				"required": ["name"]
			}`),
			wantErr: true,
		},
		{
			name:      "Empty schema",
			schemaDoc: []byte(``),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateSchema(tt.schemaDoc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
