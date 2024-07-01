package logx

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"os"
)

var Default *Helper

func init() {
	Default = NewHelper(SetUpLog("", "", "", &LogxConf{}))
}

func WithContext(ctx context.Context) *Helper {
	return &Helper{
		msgKey:  Default.msgKey,
		logger:  log.WithContext(ctx, Default.logger),
		field:   Default.field,
		sprint:  Default.sprint,
		sprintf: Default.sprintf,
	}
}

func WithField(fields ...string) *Helper {
	if len(fields)%2 != 0 {
		fields = append(fields, "")
	}
	Default.field = fields
	return Default
}

// Log Print log by level and keyvals.
func Log(level log.Level, keyvals ...interface{}) {
	kvs := []interface{}{}
	kvs = append(kvs, Default.field, keyvals)
	_ = Default.logger.Log(level, kvs...)
}

// Debug logs a message at debug level.
func Debug(a ...interface{}) {
	_ = Default.logger.Log(log.LevelDebug, Default.msgKey, Default.sprint(a...))
}

// Debugf logs a message at debug level.
func Debugf(format string, a ...interface{}) {
	_ = Default.logger.Log(log.LevelDebug, Default.msgKey, Default.sprintf(format, a...))
}

// Debugw logs a message at debug level.
func Debugw(keyvals ...interface{}) {
	_ = Default.logger.Log(log.LevelDebug, keyvals...)
}

// Info logs a message at info level.
func Info(a ...interface{}) {
	_ = Default.logger.Log(log.LevelInfo, Default.msgKey, Default.sprint(a...))
}

// Infof logs a message at info level.
func Infof(format string, a ...interface{}) {
	_ = Default.logger.Log(log.LevelInfo, Default.msgKey, Default.sprintf(format, a...))
}

// Infow logs a message at info level.
func Infow(keyvals ...interface{}) {
	_ = Default.logger.Log(log.LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func Warn(a ...interface{}) {
	_ = Default.logger.Log(log.LevelWarn, Default.msgKey, Default.sprint(a...))
}

// Warnf logs a message at warnf level.
func Warnf(format string, a ...interface{}) {
	_ = Default.logger.Log(log.LevelWarn, Default.msgKey, Default.sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func Warnw(keyvals ...interface{}) {
	_ = Default.logger.Log(log.LevelWarn, keyvals...)
}

// Error logs a message at error level.
func Error(a error) {
	kvs := []interface{}{}
	for _, v := range Default.field {
		kvs = append(kvs, v)
	}
	switch a.(type) {
	case *errors.Error:
		ee := a.(*errors.Error)
		kvs = append(kvs, "code", ee.Code, "reason", ee.Reason, "message", ee.Message, "metadata", ee.Metadata, "cause", ee.Unwrap())
	default:
		kvs = append(kvs, "error", a)
	}
	_ = Default.logger.Log(log.LevelError, kvs...)
}

// Errorf logs a message at error level.
func Errorf(format string, a ...interface{}) {
	_ = Default.logger.Log(log.LevelError, Default.msgKey, Default.sprintf(format, a...))
}

// Errorw logs a message at error level.
func Errorw(keyvals ...interface{}) {
	_ = Default.logger.Log(log.LevelError, keyvals...)
}

// Fatal logs a message at fatal level.
func Fatal(a ...interface{}) {
	_ = Default.logger.Log(log.LevelFatal, Default.msgKey, Default.sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, a ...interface{}) {
	_ = Default.logger.Log(log.LevelFatal, Default.msgKey, Default.sprintf(format, a...))
	os.Exit(1)
}

// Fatalw logs a message at fatal level.
func Fatalw(keyvals ...interface{}) {
	_ = Default.logger.Log(log.LevelFatal, keyvals...)
	os.Exit(1)
}
