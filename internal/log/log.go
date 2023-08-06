package log

var (
	lg Logger
)

// Init logger. Should be called once in application lifecycle
func Init(debug bool) {
	lg = NewStandardLogger(debug)
	// You can replace logger using broadly if it implements Logger interface
	// var l *zap.Logger
	// if debug {
	// 	l, err = zap.NewDevelopment()
	// } else {
	// 	l, err = zap.NewProduction()
	// }
	// lg = &zapLogger{
	// 	lgr: l,
	// }
}

// Debug logs the vals at Debug level.
func Debug(vals ...interface{}) {
	lg.Debug(vals...)
}

// Debugf logs the formatted message at Debug level.
func Debugf(format string, vals ...interface{}) {
	lg.Debugf(format, vals...)
}

// Info logs the vals at Info level.
func Info(vals ...interface{}) {
	lg.Info(vals...)
}

// Infof logs the formatted message at Info level.
func Infof(format string, vals ...interface{}) {
	lg.Infof(format, vals...)
}

// Error logs the vals at Error level.
func Error(vals ...interface{}) {
	lg.Error(vals...)
}

// Errorf logs the formatted message at Error level.
func Errorf(format string, vals ...interface{}) {
	lg.Errorf(format, vals...)
}

// Fatal logs the vals at Fatal level, then calls os.Exit(1).
func Fatal(vals ...interface{}) {
	lg.Fatal(vals...)
}

// Fatalf logs the formatted message at Fatal level, then calls os.Exit(1).
func Fatalf(format string, vals ...interface{}) {
	lg.Fatalf(format, vals...)
}
