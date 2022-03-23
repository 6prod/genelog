package level

import (
	"io"

	"github.com/6prod/genelog"
)

type LevelLogger struct {
	*genelog.Logger
}

func NewLevelLogger(w io.Writer) LevelLogger {
	return LevelLogger{
		genelog.New(w),
	}.AddHook(HookLevelSkip)
}

func (l LevelLogger) Info(v ...interface{}) {
	Info(l.Logger, v...)
}

func (l LevelLogger) Infoln(v ...interface{}) {
	Infoln(l.Logger, v...)
}

func (l LevelLogger) Infof(format string, v ...interface{}) {
	Infof(l.Logger, format, v...)
}

func (l LevelLogger) Error(v ...interface{}) {
	Error(l.Logger, v...)
}

func (l LevelLogger) Errorln(v ...interface{}) {
	Errorln(l.Logger, v...)
}

func (l LevelLogger) Errorf(format string, v ...interface{}) {
	Errorf(l.Logger, format, v...)
}

func (l LevelLogger) Debug(v ...interface{}) {
	Debug(l.Logger, v...)
}

func (l LevelLogger) Debugln(v ...interface{}) {
	Debugln(l.Logger, v...)
}

func (l LevelLogger) Debugf(format string, v ...interface{}) {
	Debugf(l.Logger, format, v...)
}

func (l LevelLogger) Warning(v ...interface{}) {
	Warning(l.Logger, v...)
}

func (l LevelLogger) Warningln(v ...interface{}) {
	Warningln(l.Logger, v...)
}

func (l LevelLogger) Warningf(format string, v ...interface{}) {
	Warningf(l.Logger, format, v...)
}

func (l LevelLogger) Fatal(v ...interface{}) {
	Fatal(l.Logger, v...)
}

func (l LevelLogger) Fatalln(v ...interface{}) {
	Fatalln(l.Logger, v...)
}

func (l LevelLogger) Fatalf(format string, v ...interface{}) {
	Fatalf(l.Logger, format, v...)
}

func (l LevelLogger) WithContext(v interface{}) LevelLogger {
	logger := l.Logger.WithContext(v)
	return LevelLogger{logger}
}

func (l LevelLogger) WithFormatter(f genelog.Format) LevelLogger {
	logger := l.Logger.WithFormatter(f)
	return LevelLogger{logger}
}

func (l LevelLogger) AddHook(h genelog.Hook) LevelLogger {
	logger := l.Logger.AddHook(h)
	return LevelLogger{logger}
}
