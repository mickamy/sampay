package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

const callerSkipDepth = 2

func init() {
	slog.SetDefault(slog.New(jsonHandler()))
}

func jsonHandler() *slog.JSONHandler {
	var opts = &slog.HandlerOptions{
		Level: logLevel(),
	}
	if config.Common().Env.ShouldLogToFile() {
		mw := io.MultiWriter(os.Stdout, fileWriter(config.Common().Env))
		return slog.NewJSONHandler(mw, opts)
	}
	return slog.NewJSONHandler(os.Stdout, opts)
}

func logLevel() slog.Level {
	cfg := config.Common()
	switch cfg.LogLevel {
	case config.LogLevelDebug:
		return slog.LevelDebug
	case config.LogLevelInfo:
		return slog.LevelInfo
	case config.LogLevelWarn:
		return slog.LevelWarn
	case config.LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func internalHandle(ctx context.Context, level slog.Level, msg string, args ...any) {
	_, f, l, _ := runtime.Caller(callerSkipDepth)
	source := f + ":" + strconv.Itoa(l)

	source = strings.TrimPrefix(source, config.Common().ModuleRoot+"/")

	handle(ctx, level, msg, source, args...)
}

func handle(ctx context.Context, level slog.Level, msg, source string, args ...any) {
	executionID, _ := contexts.ExecutionID(ctx)
	if executionID != uuid.Nil {
		args = append(args, slog.String("execution_id", executionID.String()))
	}

	userID, _ := contexts.AuthenticatedUserID(ctx)
	if userID != "" {
		args = append(args, slog.String("user_id", userID))
	}

	systemUserID, _ := contexts.SystemUserID(ctx)
	if systemUserID != "" {
		args = append(args, slog.String("system_user_id", systemUserID))
	}

	args = append(args, slog.String("source", source))

	slog.Default().Log(ctx, level, msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	internalHandle(ctx, slog.LevelDebug, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	internalHandle(ctx, slog.LevelInfo, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	internalHandle(ctx, slog.LevelWarn, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	internalHandle(ctx, slog.LevelError, msg, args...)
}
