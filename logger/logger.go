package logger

const (
	// LogLevelDebug is a debug logging leve.
	LogLevelDebug = "debug"
	// LogLevelInfo is an info logging leve.
	LogLevelInfo = "info"
	// LogLevelError is an error logging leve.
	LogLevelError = "error"
)

// Logger is an interface for logger which should be provided on initialization.
type Logger interface {
	Log(level, message string, additionals map[string]interface{})
	Logf(level, template string, additionals map[string]interface{}, data ...interface{})
}

// DefaultLogger is an implement Logger interface.
type DefaultLogger struct{}

// Log logs passed message.
func (d DefaultLogger) Log(_, _ string, _ map[string]interface{}) {}

// Logf logs passed message using template.
func (d DefaultLogger) Logf(_, _ string, _ map[string]interface{}, _ ...interface{}) {}
