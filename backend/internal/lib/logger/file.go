package logger

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/mickamy/sampay/config"
)

func fileWriter(env config.Env) io.Writer {
	return &lumberjack.Logger{
		Filename:   logFilePath(env),
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	}
}

func logFilePath(env config.Env) string {
	filename := filepath.Base(os.Args[0]) + ".log"
	switch env {
	case config.EnvDevelopment, config.EnvTest:
		return path.Join(os.TempDir(), filename)
	case config.EnvStaging, config.EnvProduction:
		return path.Join("var", "log", "backend-template", filename)
	default:
		return path.Join("var", "log", "backend-template", filename)
	}
}
