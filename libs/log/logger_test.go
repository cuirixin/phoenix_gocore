package log

import (
	"log"
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func TestLoggerNormal(t *testing.T) {
	Init()
	Info("abcdddfadfafafafafafafafafafafafa")
}

func TestGolangLog(t *testing.T) {
	logger := log.New(&lumberjack.Logger{
		Filename:   "/tmp/logs",
		MaxSize:    50, // megabytes
		MaxBackups: 10,
		MaxAge:     28, // days
	}, "", 0)

	logger.Println("abcdddfadfafafafafafafafafafafafa")
	teardown("/tmp/logs")
}

func TestLoggerDataPrint(t *testing.T) {

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "/tmp/data",
		MaxSize:    50, // megabytes
		MaxBackups: 10,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		w, // os.Stdout while debug
		zap.InfoLevel,
	)
	logger := zap.New(core)
	logger.Info("fff")
}

// BenchmarkZapLog-8   	  200000	      7371 ns/op
func BenchmarkZapLog(b *testing.B) {
	Init()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("abcdddfadfafafafafafafafafafafafa")
	}
	b.StopTimer()
	teardown("/tmp/logs")
}

// BenchmarkGolangLog-8   	  200000	      5849 ns/op
func BenchmarkGolangLog(b *testing.B) {
	logger := log.New(&lumberjack.Logger{
		Filename:   "/tmp/logs",
		MaxSize:    50, // megabytes
		MaxBackups: 10,
		MaxAge:     28, // days
	}, "", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Println("abcdddfadfafafafafafafafafafafafa")
	}
	b.StopTimer()
	teardown("/tmp/logs")
}

func teardown(filePath string) {
	os.Remove(filePath)
}
