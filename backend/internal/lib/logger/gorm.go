package logger

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	lib "gorm.io/gorm/logger"

	"mickamy.com/sampay/config"
)

var (
	Gorm = New(log.New(os.Stdout, "\r\n", log.LstdFlags), lib.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  lib.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})
	loggerFile string
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	loggerFile = file
}

// New initialize logger
func New(writer lib.Writer, config lib.Config) lib.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = lib.Green + "%s\n" + lib.Reset + lib.Green + "[info] " + lib.Reset
		warnStr = lib.BlueBold + "%s\n" + lib.Reset + lib.Magenta + "[warn] " + lib.Reset
		errStr = lib.Magenta + "%s\n" + lib.Reset + lib.Red + "[error] " + lib.Reset
		traceStr = lib.Green + "%s\n" + lib.Reset + lib.Yellow + "[%.3fms] " + lib.BlueBold + "[rows:%v]" + lib.Reset + " %s"
		traceWarnStr = lib.Green + "%s " + lib.Yellow + "%s\n" + lib.Reset + lib.RedBold + "[%.3fms] " + lib.Yellow + "[rows:%v]" + lib.Magenta + " %s" + lib.Reset
		traceErrStr = lib.RedBold + "%s " + lib.MagentaBold + "%s\n" + lib.Reset + lib.Yellow + "[%.3fms] " + lib.BlueBold + "[rows:%v]" + lib.Reset + " %s"
	}

	return &gormLogger{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type gormLogger struct {
	lib.Writer
	lib.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *gormLogger) LogMode(level lib.LogLevel) lib.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= lib.Info {
		l.Printf(l.infoStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= lib.Warn {
		l.Printf(l.warnStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= lib.Error {
		l.Printf(l.errStr+msg, append([]interface{}{fileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= lib.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= lib.Error && (!errors.Is(err, lib.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceErrStr, fileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceErrStr, fileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= lib.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Printf(l.traceWarnStr, fileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceWarnStr, fileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == lib.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceStr, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceStr, fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

// ParamsFilter filter params
func (l *gormLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.Config.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}

func fileWithLineNum() string {
	pcs := [13]uintptr{}
	// the third caller usually from gorm internal
	length := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:length])
	for i := 0; i < length; i++ {
		// second return value is "more", not "ok"
		frame, _ := frames.Next()
		if !isIgnoreFile(frame.File) && !strings.HasSuffix(frame.File, ".gen.go") {
			s := string(strconv.AppendInt(append([]byte(frame.File), ':'), int64(frame.Line), 10))
			return strings.TrimPrefix(s, config.Common().PackageRoot+"/")
		}
	}

	return ""
}

func isIgnoreFile(file string) bool {
	return isGormDir(file) || isLoggerFile(file) || isRepositoryFile(file)
}

func isGormDir(file string) bool {
	return strings.Contains(file, "gorm.io")
}

func isLoggerFile(file string) bool {
	return strings.HasPrefix(file, loggerFile)
}

func isRepositoryFile(file string) bool {
	return strings.HasSuffix(file, "repository.go")
}
