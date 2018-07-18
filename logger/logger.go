package logger

import (
	"fmt"
	"os"

	"bytes"
	"github.com/mattn/go-isatty"
	"github.com/op/go-logging"
)

var (
	errCantCreateLogger = fmt.Errorf("can't create logger")
	errInvalidLogger    = fmt.Errorf("invalid logger name")
)

type Level = logging.Level

// Log levels.
const (
	CRITICAL Level = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
	NONE
)

type Password string

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

type Logger struct {
	*logging.Logger
	writerLevel Level
	writer      *os.File
}

var loggers = make(map[string]*Logger)

type LoggerConfig struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Format string `yaml:"format"`
	Level  string `yaml:"level"`
	Dest   string `yaml:"dest"`
}

func SetupLoggers(loggersConfig []LoggerConfig) error {
	for _, l := range loggersConfig {
		if _, ok := loggers[l.Name]; !ok {
			var logio *os.File
			var err error

			switch l.Type {
			case "stdout":
				logio = os.Stdout
			case "stderr":
				logio = os.Stderr
			case "file":
				logio, err = os.OpenFile(l.Dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					return fmt.Errorf("can't create logger: %s", err.Error())
				}
			}
			log, err := SetupLog(l.Name, logio, l.Level, l.Format)
			if err != nil {
				return errCantCreateLogger
			}
			loggers[l.Name] = log
		}
	}
	return nil
}

func GetLogger(module string) *Logger {
	if l, ok := loggers[module]; ok {
		return l
	}
	panic(errInvalidLogger)
}

func GetLoggerErr(module string) (*Logger, error) {
	if l, ok := loggers[module]; ok {
		return l, nil
	}
	return nil, fmt.Errorf("logger `%s` not found", module)
}

func SetupLogWithWriter(module string, writer *os.File, level string, format string, writerLevel string) (*Logger, error) {

	if isatty.IsTerminal(writer.Fd()) {
		format = "%{color}" + format + "%{color:reset}"
	}
	formatter := logging.MustStringFormatter(format)
	backend := logging.NewLogBackend(writer, "", 0)
	backendWithFormater := logging.NewBackendFormatter(backend, formatter)
	leveledBackend := logging.AddModuleLevel(backendWithFormater)
	levelInt, err := logging.LogLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level")
	}

	leveledBackend.SetLevel(levelInt, module)

	log := logging.MustGetLogger(module)
	log.SetBackend(leveledBackend)

	return &Logger{log, Level(levelInt), writer}, nil
}

func SetupLog(module string, writer *os.File, level string, format string) (*Logger, error) {
	return SetupLogWithWriter(module, writer, level, format, "none")
}

func (l *Logger) D(format string, args ...interface{}) {
	l.ExtraCalldepth = 1
	l.Debugf(format, args...)
}

func (l *Logger) Dc(c int, format string, args ...interface{}) {
	l.ExtraCalldepth = 1 + c // skip `c` calls
	l.Debugf(format, args...)
}

func (l *Logger) I(format string, args ...interface{}) {
	l.ExtraCalldepth = 1
	l.Infof(format, args...)
}

func (l *Logger) Ic(c int, format string, args ...interface{}) {
	l.ExtraCalldepth = 1 + c
	l.Infof(format, args...)
}

func (l *Logger) N(format string, args ...interface{}) {
	l.ExtraCalldepth = 1
	l.Noticef(format, args...)
}

func (l *Logger) Nc(c int, format string, args ...interface{}) {
	l.ExtraCalldepth = 1 + c
	l.Noticef(format, args...)
}

func (l *Logger) W(format string, args ...interface{}) {
	l.ExtraCalldepth = 1
	l.Warningf(format, args...)
}

func (l *Logger) Wc(c int, format string, args ...interface{}) {
	l.ExtraCalldepth = 1 + c
	l.Warningf(format, args...)
}

func (l *Logger) E(format string, args ...interface{}) {
	l.ExtraCalldepth = 1
	l.Errorf(format, args...)
}

func (l *Logger) Ec(c int, format string, args ...interface{}) {
	l.ExtraCalldepth = 1 + c
	l.Errorf(format, args...)
}

func (l *Logger) C(format string, args ...interface{}) {
	l.ExtraCalldepth = 1
	l.Critical(format, args...)
}

func (l *Logger) Cc(c int, format string, args ...interface{}) {
	l.ExtraCalldepth = 1 + c
	l.Critical(format, args...)
}

func (l *Logger) Write(p []byte) (n int, err error) {
	if l.writerLevel == NONE {

		//fmt.Printf("WRITE CALLED\n")
	}

	switch l.writerLevel {
	case NONE:
		break
	case DEBUG:
		l.D(string(bytes.Trim(p, "\n")))
	case INFO:
		l.I(string(bytes.Trim(p, "\n")))
	case NOTICE:
		l.N(string(bytes.Trim(p, "\n")))
	case WARNING:
		l.W(string(bytes.Trim(p, "\n")))
	case ERROR:
		l.E(string(bytes.Trim(p, "\n")))
	case CRITICAL:
		l.C(string(bytes.Trim(p, "\n")))
	}
	return len(p), nil
}

func (l *Logger) Close() (err error) {
	return l.writer.Close()
}
