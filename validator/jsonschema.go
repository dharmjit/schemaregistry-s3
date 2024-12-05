package validator

import (
	"fmt"
	"strings"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

type JsonSchemaSpecValidator struct {
}

func NewJsonSchemaSpecValidator() *JsonSchemaSpecValidator {
	return &JsonSchemaSpecValidator{}
}

// ValidateSchema validate the json schema doc itself
func (v *JsonSchemaSpecValidator) ValidateSchema(schemaDoc []byte) error {
	c := jsonschema.NewCompiler()
	c.DefaultDraft(jsonschema.Draft2020)
	c.AssertFormat()

	doc, err := jsonschema.UnmarshalJSON(strings.NewReader(string(schemaDoc)))
	if err != nil {
		return fmt.Errorf("failed to unmarshal schema: %w", err)
	}
	err = c.AddResource("schema.json", doc)
	if err != nil {
		return fmt.Errorf("failed to add schema: %w", err)
	}

	_, err = c.Compile("schema.json")
	if err != nil {
		return fmt.Errorf("schema is invalid: %v", err)
	}
	return nil
}
