package profanity

// Logger are the methods required on the logger.
type Logger interface {
	Printf(string, ...interface{})
	Errorf(string, ...interface{})
}
