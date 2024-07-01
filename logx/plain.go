package logx

import (
	"bytes"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	golog "log"

	"sync"
)

var _ log.Logger = (*plainLogger)(nil)

type plainLogger struct {
	log  *golog.Logger
	pool *sync.Pool
}

// NewStdLogger new a logger with writer.
func newPlainLogger(w io.Writer) *plainLogger {
	return &plainLogger{
		log: golog.New(w, "", 0),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Log print the kv pairs log.
func (l *plainLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}
	if (len(keyvals) & 1) == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}
	buf := l.pool.Get().(*bytes.Buffer)

	for i := 0; i < len(keyvals); i += 2 {
		if i == 0 {
			_, _ = fmt.Fprintf(buf, "%s=%v level=%s", keyvals[i], keyvals[i+1], level.String())
			continue
		}
		_, _ = fmt.Fprintf(buf, " %s=%v", keyvals[i], keyvals[i+1])
	}
	_ = l.log.Output(4, buf.String()) //nolint:gomnd
	buf.Reset()
	l.pool.Put(buf)
	return nil
}

func (l *plainLogger) Close() error {
	return nil
}
