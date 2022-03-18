package level

import (
	"fmt"

	"github.com/6prod/genelog"
)

func Info(logger genelog.Logger, v ...interface{}) {
	Output(logger, INFO, func(logger genelog.Logger) error {
		logger.Print(v...)
		return nil
	})
}

func Infoln(logger genelog.Logger, v ...interface{}) {
	Output(logger, INFO, func(logger genelog.Logger) error {
		logger.Println(v...)
		return nil
	})
}

func Infof(logger genelog.Logger, format string, v ...interface{}) {
	Output(logger, INFO, func(logger genelog.Logger) error {
		logger.Printf(format, v...)
		return nil
	})
}

func Error(logger genelog.Logger, v ...interface{}) {
	Output(logger, ERROR, func(logger genelog.Logger) error {
		logger.Print(v...)
		return nil
	})
}

func Errorln(logger genelog.Logger, v ...interface{}) {
	Output(logger, ERROR, func(logger genelog.Logger) error {
		logger.Println(v...)
		return nil
	})
}

func Errorf(logger genelog.Logger, format string, v ...interface{}) {
	Output(logger, ERROR, func(logger genelog.Logger) error {
		logger.Printf(format, v...)
		return nil
	})
}

func Debug(logger genelog.Logger, v ...interface{}) {
	Output(logger, DEBUG, func(logger genelog.Logger) error {
		logger.Print(v...)
		return nil
	})
}

func Debugln(logger genelog.Logger, v ...interface{}) {
	Output(logger, DEBUG, func(logger genelog.Logger) error {
		logger.Println(v...)
		return nil
	})
}

func Debugf(logger genelog.Logger, format string, v ...interface{}) {
	Output(logger, DEBUG, func(logger genelog.Logger) error {
		logger.Printf(format, v...)
		return nil
	})
}

func Warning(logger genelog.Logger, v ...interface{}) {
	Output(logger, WARNING, func(logger genelog.Logger) error {
		logger.Print(v...)
		return nil
	})
}

func Warningln(logger genelog.Logger, v ...interface{}) {
	Output(logger, WARNING, func(logger genelog.Logger) error {
		logger.Println(v...)
		return nil
	})
}

func Warningf(logger genelog.Logger, format string, v ...interface{}) {
	Output(logger, WARNING, func(logger genelog.Logger) error {
		logger.Printf(format, v...)
		return nil
	})
}

func Fatal(logger genelog.Logger, v ...interface{}) {
	Output(logger, FATAL, func(logger genelog.Logger) error {
		logger.Print(v...)
		return nil
	})
}

func Fatalln(logger genelog.Logger, v ...interface{}) {
	Output(logger, FATAL, func(logger genelog.Logger) error {
		logger.Println(v...)
		return nil
	})
}

func Fatalf(logger genelog.Logger, format string, v ...interface{}) {
	Output(logger, FATAL, func(logger genelog.Logger) error {
		logger.Printf(format, v...)
		return nil
	})
}

func Output(logger genelog.Logger, level Level, output func(logger genelog.Logger) error) error {
	context, ok := logger.Context().(Leveler)
	if !ok {
		return fmt.Errorf("logger: Leveler interface not implemented")
	}

	if !IsActive(context.LevelMin(), level) {
		return nil
	}

	context.LevelSet(level)
	logger = logger.WithContext(context)

	return output(logger)
}
