package logx

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"runtime"
	"strconv"
	"strings"
)

func SetUpLog(serviceId string, serviceName string, serviceVersion string, c *LogxConf) log.Logger {
	filename := "service.log"
	if c.FileName != "" {
		filename = c.FileName
	}
	l := log.With(
		NewLogger(
			Dir(c.PathDir),
			Filename(filename),
			MaxSize(int(c.MaxSize)),
			KeepDay(int(c.KeepDays)),
			MaxBackup(int(c.MaxBackups)),
			Encoding(c.Encoding.String()),
			Level(c.Level.String()),
			FilterKey("args"),
			//logx.EntryptionFn(logEntryption),
		),
		"ts", log.DefaultTimestamp,
		//"caller", log.DefaultCaller,
		"caller", Caller(4),
		"service.id", serviceId,
		"service.name", serviceName,
		"service.version", serviceVersion,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	//return log.NewFilter(l)
	return l
}

func Caller(depth int) log.Valuer {
	return func(context.Context) interface{} {
		_, file, line, _ := runtime.Caller(depth)
		//fmt.Println(file)
		idx1 := strings.LastIndex(file, "vendor")
		if idx1 == -1 {
			idx := strings.LastIndexByte(file, '/')
			if idx == -1 {
				return file[idx+1:] + ":" + strconv.Itoa(line)
			}
			idx = strings.LastIndexByte(file[:idx], '/')
			idx = strings.LastIndexByte(file[:idx], '/')
			idx = strings.LastIndexByte(file[:idx], '/')
			return file[idx+1:] + ":" + strconv.Itoa(line)
		}
		return strings.TrimPrefix(file[idx1:], "vendor/") + ":" + strconv.Itoa(line)
	}
}
