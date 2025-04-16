package log

import (
	"context"
	"fmt"
	"github.com/Deeksharma/taskmanager/internal/config"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

var entry *log.Entry

func init() {
	var logger *log.Logger = log.StandardLogger()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.SetLevel(log.Level(config.GetInt32("logging.level")))
	entry = log.NewEntry(logger)

	//file, err := os.OpenFile(filepath.Join(config.GetString("logging.dir_name"), config.GetString("logging.file_name")), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	//if err == nil {
	// logger.Out = file
	//} else {
	// logger.Info("Failed to log to file, using default stderr, err: ", err)
	//}

}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 0
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func Info(ctx context.Context, msg string) {
	entry.Data["caller"] = fileInfo(2)
	entry.Info(msg)
}

func Infof(ctx context.Context, msg string, args ...interface{}) {
	entry.Data["caller"] = fileInfo(2)
	entry.Infof(msg, args...)
}

func Fatal(ctx context.Context, msg string) {
	entry.Data["caller"] = fileInfo(2)
	entry.Fatal(msg)
}

func Fatalf(ctx context.Context, msg string, args ...interface{}) {
	entry.Data["caller"] = fileInfo(2)
	entry.Fatalf(msg, args...)
}

func Error(ctx context.Context, msg string) {
	entry.Data["caller"] = fileInfo(2)
	entry.Error(msg)
}

func Errorf(ctx context.Context, msg string, args ...interface{}) {
	entry.Data["caller"] = fileInfo(2)
	entry.Errorf(msg, args...)
}

func Debug(ctx context.Context, msg string) {
	entry.Data["caller"] = fileInfo(2)
	entry.Debug(msg)
}

func Debugf(ctx context.Context, msg string) {
	entry.Data["caller"] = fileInfo(2)
	entry.Debugf(msg)
}

func InfoWithFields(ctx context.Context, fields map[string]interface{}, msg string) {
	entry.Data["caller"] = fileInfo(2)
	entry.WithFields(fields).Info(msg)
}

func ErrorWithFields(ctx context.Context, fields map[string]interface{}, msg string) {
	entry.Data["caller"] = fileInfo(2)
	entry.WithFields(fields).Error(msg)
}

func DebugWithFields(ctx context.Context, fields map[string]interface{}, msg string) {
	entry.Data["caller"] = fileInfo(2)
	entry.WithFields(fields).Debug(msg)
}
