package log

import "go.uber.org/zap"

// zapLogger implements Logger interface
type zapLogger struct {
	lgr *zap.Logger
}

// Debug logs the vals at Debug level.
func (zl *zapLogger) Debug(vals ...interface{}) {
	zl.lgr.Sugar().Debug(vals...)
}

// Debugf logs the formatted message at Debug level.
func (zl *zapLogger) Debugf(format string, vals ...interface{}) {
	zl.lgr.Sugar().Debugf(format, vals...)
}

// Info logs the vals at Info level.
func (zl *zapLogger) Info(vals ...interface{}) {
	zl.lgr.Sugar().Info(vals...)
}

// Infof logs the formatted message at Info level.
func (zl *zapLogger) Infof(format string, vals ...interface{}) {
	zl.lgr.Sugar().Infof(format, vals...)
}

// Error logs the vals at Error level.
func (zl *zapLogger) Error(vals ...interface{}) {
	zl.lgr.Sugar().Error(vals...)
}

// Errorf logs the formatted message at Error level.
func (zl *zapLogger) Errorf(format string, vals ...interface{}) {
	zl.lgr.Sugar().Errorf(format, vals...)
}

// Fatal logs the vals at Fatal level, then calls os.Exit(1).
func (zl *zapLogger) Fatal(vals ...interface{}) {
	zl.lgr.Sugar().Fatal(vals...)
}

// Fatalf logs the formatted message at Fatal level, then calls os.Exit(1).
func (zl *zapLogger) Fatalf(format string, vals ...interface{}) {
	zl.lgr.Sugar().Fatalf(format, vals...)
}

// Close calls finalizer of logger implementations.
func (zl *zapLogger) Close() error {
	return zl.lgr.Sync()
}
