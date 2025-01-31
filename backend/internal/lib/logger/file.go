package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

func FileWriter() (io.Writer, error) {
	filename := filepath.Base(os.Args[0]) + ".log"
	p := path.Join("/var", "log", "sampay", filename)

	dir := path.Dir(p)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	return &lumberjack.Logger{
		Filename:   p,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}, nil
}
