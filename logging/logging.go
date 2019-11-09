package logging

import (
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

type ContextHook struct {
}

func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
func (hook ContextHook) Fire(entry *logrus.Entry) error {
	if pc, file, line, ok := runtime.Caller(7); ok {
		funcName := runtime.FuncForPC(pc).Name()
		entry.Data["file"] = path.Base(file)
		entry.Data["func"] = path.Base(funcName)
		entry.Data["line"] = line
	}
	return nil
}

func GetLogger() *Logger {
	l := logrus.New()

	l.SetLevel(logrus.DebugLevel)
	l.AddHook(ContextHook{})
	l.Formatter = &easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "%time% [%lvl%] %func%(%line%) %msg%\n",
	}
	return (*Logger)(l)
}

// alias for logrus.Logger,
// or for logrus.Entry, if uses WithFields in GetLogger .
type Logger = logrus.Logger
