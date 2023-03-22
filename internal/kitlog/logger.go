package kitlog

type (
	Logger struct {
		log logger
	}

	logger interface {
		Println(args ...interface{})
	}
)

// NewLogger returns a new Logger that wraps the provided logger.
func NewLogger(l logger) Logger {
	return Logger{log: l}
}

// Log implements the Log method of the go-kit log.Logger interface.
func (l Logger) Log(keyvals ...interface{}) error {
	l.log.Println(keyvals...)
	return nil
}
