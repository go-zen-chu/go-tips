package log

import (
	"log"
	"os"
)

// sample implementation using golang standard logger
type standardLogger struct {
	logger log.Logger
	debug  bool
}

func NewStandardLogger(debug bool) *standardLogger {
	sl := &standardLogger{
		logger: *log.Default(),
		debug:  debug,
	}
	return sl
}

// Debug logs the vals at Debug level.
func (s *standardLogger) Debug(vals ...interface{}) {
	if s.debug {
		s.logger.Print(vals...)
	}
}

// Debugf logs the formatted message at Debug level.
func (s *standardLogger) Debugf(format string, vals ...interface{}) {
	if s.debug {
		s.logger.Printf(format, vals...)
	}
}

// Info logs the vals at Info level.
func (s *standardLogger) Info(vals ...interface{}) {
	s.logger.Print(vals...)
}

// Infof logs the formatted message at Info level.
func (s *standardLogger) Infof(format string, vals ...interface{}) {
	s.logger.Printf(format, vals...)
}

// Error logs the vals at Error level.
func (s *standardLogger) Error(vals ...interface{}) {
	s.logger.SetOutput(os.Stderr)
	s.logger.Print(vals...)
	s.logger.SetOutput(os.Stdout)
}

// Errorf logs the formatted message at Error level.
func (s *standardLogger) Errorf(format string, vals ...interface{}) {
	s.logger.SetOutput(os.Stderr)
	s.logger.Printf(format, vals...)
	s.logger.SetOutput(os.Stdout)
}

// Fatal logs the vals at Fatal level, then calls os.Exit(1).
func (s *standardLogger) Fatal(vals ...interface{}) {
	s.logger.SetOutput(os.Stderr)
	s.logger.Fatal(vals...)
	s.logger.SetOutput(os.Stdout)
}

// Fatalf logs the formatted message at Fatal level, then calls os.Exit(1).
func (s *standardLogger) Fatalf(format string, vals ...interface{}) {
	s.logger.SetOutput(os.Stderr)
	s.logger.Fatalf(format, vals...)
	s.logger.SetOutput(os.Stdout)
}
