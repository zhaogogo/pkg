package logx

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/rs/zerolog"
	"io"
)

type ZeroLog struct {
	log zerolog.Logger
}

func NewZeroLoggerx(w io.Writer) *ZeroLog {
	l := zerolog.New(w)
	//switch level {
	//case log.LevelDebug.String():
	//	l = l.Level(zerolog.DebugLevel)
	//case log.LevelInfo.String():
	//	l = l.Level(zerolog.InfoLevel)
	//case log.LevelWarn.String():
	//	l = l.Level(zerolog.WarnLevel)
	//case log.LevelError.String():
	//	l = l.Level(zerolog.ErrorLevel)
	//case log.LevelFatal.String():
	//	l = l.Level(zerolog.FatalLevel)
	//default:
	//	l = l.Level(zerolog.InfoLevel)
	//}
	return &ZeroLog{log: l}
}

func (l *ZeroLog) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn().Msgf("Keyvalues must appear in pairs: ", keyvals)
		return nil
	}
	var e *zerolog.Event
	switch level {
	case log.LevelDebug:
		e = l.log.Debug()
		if !e.Enabled() {
			return nil
		}
	case log.LevelInfo:
		e = l.log.Info()
		if !e.Enabled() {
			return nil
		}
	case log.LevelWarn:
		e = l.log.Warn()
		if !e.Enabled() {
			return nil
		}
	case log.LevelError:
		e = l.log.Error()
		if !e.Enabled() {
			return nil
		}
	default:
		e = l.log.Info()
		if !e.Enabled() {
			return nil
		}
	}

	msg := ""
	for i := 0; i < len(keyvals); i += 2 {
		vi := i + 1
		if fmt.Sprint(keyvals[i]) == "msg" {
			msg = fmt.Sprint(keyvals[vi])
			continue
		}

		e = e.Any(fmt.Sprint(keyvals[i]), keyvals[vi])
	}

	if msg == "" {
		e.Send()
	} else {
		e.Msg(msg)
	}

	return nil
}
