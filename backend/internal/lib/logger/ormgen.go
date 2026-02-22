package logger

import (
	"context"
	"log/slog"
	"runtime"
	"strconv"
	"strings"

	"github.com/mickamy/sampay/config"
)

var ormSkipPrefixes = []string{
	"github.com/mickamy/sampay/internal/lib/logger",
	"github.com/mickamy/ormgen/",
}

type ORM struct{}

func (l ORM) Log(ctx context.Context, query string, args ...any) {
	source := ormCallerSource()
	handle(ctx, slog.LevelDebug, "ormgen", source, "query", query, "args", args)
}

func ormCallerSource() string {
	var pcs [16]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	moduleRoot := config.Common().ModuleRoot + "/"

	for {
		frame, more := frames.Next()
		if !shouldSkip(frame.Function, frame.File) {
			return strings.TrimPrefix(frame.File+":"+strconv.Itoa(frame.Line), moduleRoot)
		}
		if !more {
			break
		}
	}

	return ""
}

func shouldSkip(funcName, fileName string) bool {
	for _, prefix := range ormSkipPrefixes {
		if strings.Contains(funcName, prefix) {
			return true
		}
	}
	return isRepositoryFrame(fileName) || isQueryFrame(fileName)
}

func isRepositoryFrame(fileName string) bool {
	return strings.HasSuffix(fileName, "_repository.go")
}

func isQueryFrame(fileName string) bool {
	return strings.HasSuffix(fileName, "_query_gen.go")
}
