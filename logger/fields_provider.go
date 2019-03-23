package logger

// FieldsProvider is a provider for fields.
type FieldsProvider interface {
	Fields() Fields
}
