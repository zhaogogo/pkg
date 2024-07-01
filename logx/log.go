package logx

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/go-kratos/kratos/v2/log"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	encoding_json  = "json"
	encoding_plain = "plain"
)

type Option func(l *Logger)

func Dir(dir string) Option {
	return func(l *Logger) { l.dir = dir }
}

func MaxSize(maxSize int) Option {
	return func(l *Logger) {
		l.maxSize = maxSize
	}
}

func KeepDay(keepDay int) Option {
	return func(l *Logger) {
		l.keepDay = keepDay
	}
}

func MaxBackup(maxBackup int) Option {
	return func(l *Logger) {
		l.maxBackups = maxBackup
	}
}

func Encoding(encoding string) Option {
	return func(l *Logger) {
		l.encoding = encoding
	}
}
func Level(level string) Option {
	return func(l *Logger) {
		l.level = log.ParseLevel(level)
	}
}

func FilterKey(keys ...string) Option {
	return func(l *Logger) {
		for _, key := range keys {
			l.filterKey[key] = struct{}{}
		}
	}
}

func EntryptionFn(fn func(string) string) Option {
	return func(l *Logger) {
		l.entryptionFn = fn
	}
}

func Filename(filename string) Option {
	return func(l *Logger) {
		l.filename = filename
	}
}

type Logger struct {
	logger log.Logger
	// 日志格式字段加密
	filterKey    map[interface{}]struct{}
	entryptionFn func(string) string

	dir        string
	filename   string
	maxSize    int
	keepDay    int
	maxBackups int
	encoding   string
	level      log.Level
}

func NewLogger(opt ...Option) *Logger {
	var (
		w io.Writer
		//accessWriter  io.Writer
	)
	logger := &Logger{
		filterKey: make(map[interface{}]struct{}),
	}

	for _, o := range opt {
		o(logger)
	}

	switch {
	case len(logger.dir) > 0:
		w = &lumberjack.Logger{
			Filename:   path.Join(logger.dir, logger.filename),
			MaxSize:    int(logger.maxSize),
			MaxAge:     int(logger.keepDay),
			MaxBackups: int(logger.maxBackups),
			Compress:   false,
		}
	case len(logger.dir) == 0:
		w = os.Stdout
	default:
		w = os.Stdout
	}

	switch logger.encoding {
	case encoding_json:
		logger.logger = NewZeroLoggerx(w)
	case encoding_plain:
		logger.logger = newPlainLogger(w)
	default:
		logger.logger = newPlainLogger(w)
	}

	return logger
}

func (lx *Logger) Log(level log.Level, keyvals ...interface{}) error {
	if lx.level > level {
		return nil
	}
	var ll log.Logger = lx.logger
	if len(lx.filterKey) > 0 && lx.entryptionFn != nil {
		for i := 0; i < len(keyvals); i += 2 {
			v := i + 1
			if v >= len(keyvals) {
				continue
			}
			if _, ok := lx.filterKey[keyvals[i]]; ok {
				keyvals[v] = lx.entryptionFn(fmt.Sprintf("%v", keyvals[v]))
			}
		}
	}
	return ll.Log(level, keyvals...)
}
