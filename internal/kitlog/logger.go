package kitlog

type (
	Logger struct {
		log logger
	}

	logger interface {
		Println(args ...interface{})
		Errorf(format string, args ...interface{})
	}
)

// NewLogger returns a new Logger that wraps the provided logger.
func NewLogger(l logger) Logger {
	return Logger{log: l}
}

// Log implements the Log method of the go-kit log.Logger interface.
func (l Logger) Log(keyvals ...interface{}) error {
	if len(keyvals) == 2 && keyvals[0] == "err" {
		l.log.Errorf("%v", keyvals[1])
		return nil
	}

	l.log.Println(keyvals...)
	return nil
}
