package log

// Logger interface is for logging interface abstraction. Any logger may satisfy this interface
type Logger interface {
	// Debug logs the vals at Debug level.
	Debug(vals ...interface{})
	// Debugf logs the formatted message at Debug level.
	Debugf(format string, vals ...interface{})
	// Info logs the vals at Info level.
	Info(vals ...interface{})
	// Infof logs the formatted message at Info level.
	Infof(format string, vals ...interface{})
	// Error logs the vals at Error level.
	Error(vals ...interface{})
	// Errorf logs the formatted message at Error level.
	Errorf(format string, vals ...interface{})
	// Fatal logs the vals at Fatal level, then calls os.Exit(1).
	Fatal(vals ...interface{})
	// Fatalf logs the formatted message at Fatal level, then calls os.Exit(1).
	Fatalf(format string, vals ...interface{})
}
