package kitlog

import "github.com/sirupsen/logrus"

type Logger struct {
	log *logrus.Entry
}

func NewLogger(l *logrus.Entry) Logger {
	return Logger{l}
}

func (l Logger) Log(keyvals ...interface{}) error {
	l.log.Println(keyvals...)
	return nil
}
