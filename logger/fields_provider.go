package logger

// FieldsProvider is a type that returns fields.
type FieldsProvider interface {
	Fields() map[string]string
}
