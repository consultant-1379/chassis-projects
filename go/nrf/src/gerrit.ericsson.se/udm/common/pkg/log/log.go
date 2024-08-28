package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"runtime"
)

//type Level uint32
type Level logrus.Level

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

type LogStruct struct {
	RequestDescription  string
	ResponseDescription string
	SequenceId          string
}

func ParseLevel(lvl string) (Level, error) {
	rlvl, err := logrus.ParseLevel(lvl)
	return Level(rlvl), err
}

func GetLevel() Level {
	return Level(logrus.GetLevel())
}

func SetLevel(level Level) {
	logrus.SetLevel(logrus.Level(level))
}

func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

func SetFormatter(formatter *JSONFormatter) {
	logrus.SetFormatter(formatter)
}

func Debug(args ...interface{}) {
	if logrus.GetLevel() >= logrus.DebugLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Debug(args...)
		} else {
			logrus.Debug(args...)
		}
	}
}

func Print(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok == true {
		logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Print(args...)
	} else {
		logrus.Print(args...)
	}
}

func Info(args ...interface{}) {
	if logrus.GetLevel() >= logrus.InfoLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Info(args...)
		} else {
			logrus.Info(args...)
		}
	}
}

func Warn(args ...interface{}) {
	if logrus.GetLevel() >= logrus.WarnLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Warn(args...)
		} else {
			logrus.Warn(args...)
		}
	}
}

func Warning(args ...interface{}) {
	if logrus.GetLevel() >= logrus.WarnLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Warning(args...)
		} else {
			logrus.Warning(args...)
		}
	}
}

func Error(args ...interface{}) {
	if logrus.GetLevel() >= logrus.ErrorLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Error(args...)
		} else {
			logrus.Error(args...)
		}
	}
}

func Panic(args ...interface{}) {
	if logrus.GetLevel() >= logrus.PanicLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Panic(args...)
		} else {
			logrus.Panic(args...)
		}
	}
}

func Fatal(args ...interface{}) {
	if logrus.GetLevel() >= logrus.FatalLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Fatal(args...)
		} else {
			logrus.Fatal(args...)
		}
	}
}

func DebugJ(args ...interface{}) {
	if logrus.GetLevel() >= logrus.DebugLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"isjson": "true", "stackinfo": fmt.Sprintf("%v:%v", file, line)}).Debug(args...)
		} else {
			logrus.Debug(args...)
		}
	}
}

func PrintJ(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok == true {
		logrus.WithFields(logrus.Fields{"isjson": "true", "stackinfo": fmt.Sprintf("%v:%v", file, line)}).Print(args...)
	} else {
		logrus.Print(args...)
	}
}

func InfoJ(args ...interface{}) {
	if logrus.GetLevel() >= logrus.InfoLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"isjson": "true", "stackinfo": fmt.Sprintf("%v:%v", file, line)}).Info(args...)
		} else {
			logrus.Info(args...)
		}
	}
}

func WarnJ(args ...interface{}) {
	if logrus.GetLevel() >= logrus.WarnLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"isjson": "true", "stackinfo": fmt.Sprintf("%v:%v", file, line)}).Warn(args...)
		} else {
			logrus.Warn(args...)
		}
	}
}

func WarningJ(args ...interface{}) {
	if logrus.GetLevel() >= logrus.WarnLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"isjson": "true", "stackinfo": fmt.Sprintf("%v:%v", file, line)}).Warning(args...)
		} else {
			logrus.Warning(args...)
		}
	}
}

func ErrorJ(args ...interface{}) {
	if logrus.GetLevel() >= logrus.ErrorLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"isjson": "true", "stackinfo": fmt.Sprintf("%v:%v", file, line)}).Error(args...)
		} else {
			logrus.Error(args...)
		}
	}
}

func PanicJ(args ...interface{}) {
	if logrus.GetLevel() >= logrus.PanicLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"isjson": "true", "stackinfo": fmt.Sprintf("%v:%v", file, line)}).Panic(args...)
		} else {
			logrus.Panic(args...)
		}
	}
}

func FatalJ(args ...interface{}) {
	if logrus.GetLevel() >= logrus.FatalLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"isjson": "true", "stackinfo": fmt.Sprintf("%v:%v", file, line)}).Fatal(args...)
		} else {
			logrus.Fatal(args...)
		}
	}
}

func Debugf(format string, args ...interface{}) {
	if logrus.GetLevel() >= logrus.DebugLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Debugf(format, args...)
		} else {
			logrus.Debugf(format, args...)
		}
	}
}

func Printf(format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok == true {
		logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Printf(format, args...)
	} else {
		logrus.Printf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if logrus.GetLevel() >= logrus.InfoLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Infof(format, args...)
		} else {
			logrus.Infof(format, args...)
		}
	}
}

func Warnf(format string, args ...interface{}) {
	if logrus.GetLevel() >= logrus.WarnLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Warnf(format, args...)
		} else {
			logrus.Warnf(format, args...)
		}
	}
}

func Warningf(format string, args ...interface{}) {
	if logrus.GetLevel() >= logrus.WarnLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Warningf(format, args...)
		} else {
			logrus.Warningf(format, args...)
		}
	}
}

func Errorf(format string, args ...interface{}) {
	if logrus.GetLevel() >= logrus.ErrorLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Errorf(format, args...)
		} else {
			logrus.Errorf(format, args...)
		}
	}
}

func Panicf(format string, args ...interface{}) {
	if logrus.GetLevel() >= logrus.PanicLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Panicf(format, args...)
		} else {
			logrus.Panicf(format, args...)
		}
	}
}

func Fatalf(format string, args ...interface{}) {
	if logrus.GetLevel() >= logrus.FatalLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Fatalf(format, args...)
		} else {
			logrus.Fatalf(format, args...)
		}
	}
}

func Debugln(args ...interface{}) {
	if logrus.GetLevel() >= logrus.DebugLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Debugln(args...)
		} else {
			logrus.Debugln(args...)
		}
	}
}

func Println(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok == true {
		logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Println(args...)
	} else {
		logrus.Println(args...)
	}
}

func Infoln(args ...interface{}) {
	if logrus.GetLevel() >= logrus.InfoLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Infoln(args...)
		} else {
			logrus.Infoln(args...)
		}
	}
}

func Warnln(args ...interface{}) {
	if logrus.GetLevel() >= logrus.WarnLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Warnln(args...)
		} else {
			logrus.Warnln(args...)
		}
	}
}

func Warningln(args ...interface{}) {
	if logrus.GetLevel() >= logrus.WarnLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Warningln(args...)
		} else {
			logrus.Warningln(args...)
		}
	}
}

func Errorln(args ...interface{}) {
	if logrus.GetLevel() >= logrus.ErrorLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Errorln(args...)
		} else {
			logrus.Errorln(args...)
		}
	}
}

func Panicln(args ...interface{}) {
	if logrus.GetLevel() >= logrus.PanicLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Panicln(args...)
		} else {
			logrus.Panicln(args...)
		}
	}
}

func Fatalln(args ...interface{}) {
	if logrus.GetLevel() >= logrus.FatalLevel {
		_, file, line, ok := runtime.Caller(1)
		if ok == true {
			logrus.WithFields(logrus.Fields{"stackinfo": fmt.Sprintf("%v:%v", file, line)}).Fatalln(args...)
		} else {
			logrus.Fatalln(args...)
		}
	}
}
