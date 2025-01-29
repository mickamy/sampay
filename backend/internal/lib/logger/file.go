package logger

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

func FileWriter() io.Writer {
	filename := filepath.Base(os.Args[0]) + ".log"
	p := path.Join(path.Dir("/"), "var", "log", "sampay", filename)
	return &lumberjack.Logger{
		Filename:   p,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}
}
