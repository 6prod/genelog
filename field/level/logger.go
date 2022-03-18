package level

import "github.com/6prod/genelog"

type LevelLogger interface {
	Info(v ...interface{})
	Infoln(v ...interface{})
	Infof(format string, v ...interface{})

	Error(v ...interface{})
	Errorln(v ...interface{})
	Errorf(format string, v ...interface{})

	Debug(v ...interface{})
	Debugln(v ...interface{})
	Debugf(format string, v ...interface{})

	Warning(v ...interface{})
	Warningln(v ...interface{})
	Warningf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalln(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type loggerWithLevel struct {
	genelog.Logger
}

func NewLevelLogger(logger genelog.Logger) loggerWithLevel {
	return loggerWithLevel{
		logger,
	}
}

func (l loggerWithLevel) Info(v ...interface{}) {
	Info(l, v...)
}

func (l loggerWithLevel) Infoln(v ...interface{}) {
	Infoln(l, v...)
}

func (l loggerWithLevel) Infof(format string, v ...interface{}) {
	Infof(l, format, v...)
}

func (l loggerWithLevel) Error(v ...interface{}) {
	Error(l, v...)
}

func (l loggerWithLevel) Errorln(v ...interface{}) {
	Errorln(l, v...)
}

func (l loggerWithLevel) Errorf(format string, v ...interface{}) {
	Errorf(l, format, v...)
}

func (l loggerWithLevel) Debug(v ...interface{}) {
	Debug(l, v...)
}

func (l loggerWithLevel) Debugln(v ...interface{}) {
	Debugln(l, v...)
}

func (l loggerWithLevel) Debugf(format string, v ...interface{}) {
	Debugf(l, format, v...)
}

func (l loggerWithLevel) Warning(v ...interface{}) {
	Warning(l, v...)
}

func (l loggerWithLevel) Warningln(v ...interface{}) {
	Warningln(l, v...)
}

func (l loggerWithLevel) Warningf(format string, v ...interface{}) {
	Warningf(l, format, v...)
}

func (l loggerWithLevel) Fatal(v ...interface{}) {
	Fatal(l, v...)
}

func (l loggerWithLevel) Fatalln(v ...interface{}) {
	Fatalln(l, v...)
}

func (l loggerWithLevel) Fatalf(format string, v ...interface{}) {
	Fatalf(l, format, v...)
}
