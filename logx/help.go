package logx

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"os"
)

// DefaultMessageKey default message key.
var DefaultMessageKey = "msg"

// HelpOption is Helper option.
type HelpOption func(*Helper)

// Helper is a logger helper.
type Helper struct {
	logger  log.Logger
	msgKey  string
	field   []string
	sprint  func(...interface{}) string
	sprintf func(format string, a ...interface{}) string
}

// WithMessageKey with message key.
func WithMessageKey(k string) HelpOption {
	return func(opts *Helper) {
		opts.msgKey = k
	}
}

// WithSprint with sprint
func WithSprint(sprint func(...interface{}) string) HelpOption {
	return func(opts *Helper) {
		opts.sprint = sprint
	}
}

// WithSprintf with sprintf
func WithSprintf(sprintf func(format string, a ...interface{}) string) HelpOption {
	return func(opts *Helper) {
		opts.sprintf = sprintf
	}
}

func WithField(fields ...string) HelpOption {
	if len(fields)%2 != 0 {
		fields = append(fields, "")
	}
	return func(opts *Helper) {
		opts.field = append(opts.field, fields...)
	}
}

// NewHelper new a logger helper.
func NewHelper(logger log.Logger, opts ...HelpOption) *Helper {
	options := &Helper{
		msgKey:  DefaultMessageKey, // default message key
		logger:  logger,
		sprint:  fmt.Sprint,
		sprintf: fmt.Sprintf,
	}
	for _, o := range opts {
		o(options)
	}
	return options
}

// WithContext returns a shallow copy of h with its context changed
// to ctx. The provided ctx must be non-nil.
func (h *Helper) WithContext(ctx context.Context) *Helper {
	return &Helper{
		msgKey:  h.msgKey,
		logger:  log.WithContext(ctx, h.logger),
		field:   h.field,
		sprint:  h.sprint,
		sprintf: h.sprintf,
	}
}

func (h *Helper) WithField(fields ...string) *Helper {
	if len(fields)%2 != 0 {
		fields = append(fields, "")
	}
	h.field = fields
	return h
}

// Log Print log by level and keyvals.
func (h *Helper) Log(level log.Level, keyvals ...interface{}) {
	kvs := []interface{}{}
	kvs = append(kvs, h.field, keyvals)
	_ = h.logger.Log(level, kvs...)
}

// Debug logs a message at debug level.
func (h *Helper) Debug(a ...interface{}) {
	_ = h.logger.Log(log.LevelDebug, h.msgKey, h.sprint(a...))
}

// Debugf logs a message at debug level.
func (h *Helper) Debugf(format string, a ...interface{}) {
	_ = h.logger.Log(log.LevelDebug, h.msgKey, h.sprintf(format, a...))
}

// Debugw logs a message at debug level.
func (h *Helper) Debugw(keyvals ...interface{}) {
	_ = h.logger.Log(log.LevelDebug, keyvals...)
}

// Info logs a message at info level.
func (h *Helper) Info(a ...interface{}) {
	_ = h.logger.Log(log.LevelInfo, h.msgKey, h.sprint(a...))
}

// Infof logs a message at info level.
func (h *Helper) Infof(format string, a ...interface{}) {
	_ = h.logger.Log(log.LevelInfo, h.msgKey, h.sprintf(format, a...))
}

// Infow logs a message at info level.
func (h *Helper) Infow(keyvals ...interface{}) {
	_ = h.logger.Log(log.LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func (h *Helper) Warn(a ...interface{}) {
	_ = h.logger.Log(log.LevelWarn, h.msgKey, h.sprint(a...))
}

// Warnf logs a message at warnf level.
func (h *Helper) Warnf(format string, a ...interface{}) {
	_ = h.logger.Log(log.LevelWarn, h.msgKey, h.sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func (h *Helper) Warnw(keyvals ...interface{}) {
	_ = h.logger.Log(log.LevelWarn, keyvals...)
}

// Error logs a message at error level.
func (h *Helper) Error(a error) {
	kvs := []interface{}{}
	for _, v := range h.field {
		kvs = append(kvs, v)
	}
	switch a.(type) {
	case *errors.Error:
		ee := a.(*errors.Error)
		kvs = append(kvs, "code", ee.Code, "reason", ee.Reason, "message", ee.Message, "metadata", ee.Metadata, "cause", ee.Unwrap())
	default:
		kvs = append(kvs, "error", a)
	}
	_ = h.logger.Log(log.LevelError, kvs...)
}

// Errorf logs a message at error level.
func (h *Helper) Errorf(format string, a ...interface{}) {
	_ = h.logger.Log(log.LevelError, h.msgKey, h.sprintf(format, a...))
}

// Errorw logs a message at error level.
func (h *Helper) Errorw(keyvals ...interface{}) {
	_ = h.logger.Log(log.LevelError, keyvals...)
}

// Fatal logs a message at fatal level.
func (h *Helper) Fatal(a ...interface{}) {
	_ = h.logger.Log(log.LevelFatal, h.msgKey, h.sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func (h *Helper) Fatalf(format string, a ...interface{}) {
	_ = h.logger.Log(log.LevelFatal, h.msgKey, h.sprintf(format, a...))
	os.Exit(1)
}

// Fatalw logs a message at fatal level.
func (h *Helper) Fatalw(keyvals ...interface{}) {
	_ = h.logger.Log(log.LevelFatal, keyvals...)
	os.Exit(1)
}
