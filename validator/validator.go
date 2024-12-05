package validator

// Validator is the interface to validate the schema
//
//go:generate go run go.uber.org/mock/mockgen -package mockschemavalidator -typed=true -source=validator.go -destination ./mocks/validator.go
type Validator interface {
	// Validate validates the schema
	ValidateSchema(schemaDoc []byte) error
}
