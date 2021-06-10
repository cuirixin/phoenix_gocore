/*
 * Copyright (c) 2018. Alibaba Cloud, All right reserved.
 * This software is the confidential and proprietary information of Alibaba Cloud ("Confidential Information").
 * You shall not disclose such Confidential Information and shall use it only in accordance with the terms of
 * the license agreement you entered into with Alibaba Cloud.
 */

package log

import (
	"encoding/json"
	"fmt"
	"github.com/cuirixin/phoenix_gocore/utils"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

type Conf struct {
	Scan       bool   `json:"scan,omitempty"`
	ScanPeriod int    `json:"scan_period,omitempty"`
	Production bool   `json:"production,omitempty"`
	Level      string `json:"level,omitempty"`
	Filename   string `json:"filename"`
	Compress   bool   `json:"compress"` // compress using gzip?
	MaxSizeMB  int    `json:"max_size_mb"`
	MaxBackups int    `json:"max_backups"`
	MaxAgeDays int    `json:"max_age_days"`
}

type Field = zap.Field

var zapLogger *zap.Logger

func InitWithConf(logConf *Conf) {
	level := level(logConf)
	encoderConfig := newEncoderConfig(logConf.Production)

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logConf.Filename,
		MaxSize:    logConf.MaxSizeMB, // megabytes
		MaxBackups: logConf.MaxBackups,
		MaxAge:     logConf.MaxAgeDays, // days
		Compress:   logConf.Compress,
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w, // os.Stdout while debug
		level,
	)
	logger := zap.New(core)

	zapLogger = logger
}

func Init() {
	logConf, err := unmarshalConfEnv()
	if err != nil {
		filePath := path.Join(utils.CallerSourcePath(), "./configs", "log.json")
		fmt.Println(filePath)
		logConf, err = unmarshalConfFile(filePath)
		if err != nil {
			logConf = &Conf{
				Scan:       false,
				ScanPeriod: 99999999,
				Production: false,
				Level:      "INFO",
				Filename:   "/tmp/phoenix-default-conf.log",
				Compress:   false,
				MaxSizeMB:  100,
				MaxBackups: 10,
				MaxAgeDays: 7,
			}
			fmt.Printf("[WARNING] initialize log conf failed, use default conf: %v\n", logConf)
		}
	}

	level := level(logConf)
	encoderConfig := newEncoderConfig(logConf.Production)

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logConf.Filename,
		MaxSize:    logConf.MaxSizeMB, // megabytes
		MaxBackups: logConf.MaxBackups,
		MaxAge:     logConf.MaxAgeDays, // days
		Compress:   logConf.Compress,
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w, // os.Stdout while debug
		level,
	)
	logger := zap.New(core)

	zapLogger = logger
}

func newEncoderConfig(production bool) zapcore.EncoderConfig {
	if production {
		return zap.NewProductionEncoderConfig()
	} else {
		return zap.NewDevelopmentEncoderConfig()
	}
}

func level(logConf *Conf) zapcore.Level {
	level := zap.InfoLevel
	switch strings.ToUpper(logConf.Level) {
	case LevelDebug:
		level = zap.DebugLevel
	case LevelInfo:
		level = zap.InfoLevel
	case LevelWarn:
		level = zap.WarnLevel
	case LevelError:
		level = zap.ErrorLevel
	}
	return level
}

func unmarshalConfFile(filename string) (*Conf, error) {
	logConf := &Conf{}

	if fi, err := os.Stat(filename); os.IsNotExist(err) || fi.IsDir() {
		return nil, fmt.Errorf("log file not found: %s\n", filename)
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read log file failed: %s, err: %s\n", filename, err.Error())
	}

	err = json.Unmarshal(content, logConf)
	if err != nil {
		return nil, fmt.Errorf("unmarshal log file failed: %s, err: %s\n", filename, err.Error())
	}

	return logConf, nil
}

func unmarshalConfEnv() (*Conf, error) {
	logKey := getComponentName() + "_LOG_CONF"
	logConfContent := os.Getenv(logKey)
	logConf := &Conf{}
	err := json.Unmarshal([]byte(logConfContent), logConf)
	if err != nil {
		return nil, fmt.Errorf("unmarshal log from env failed: %s, err: %s\n", logConfContent, err.Error())
	}

	return logConf, nil
}

func getComponentName() string {
	compName := os.Getenv("PHOENIX_COMPONENT_NAME")
	return compName
}

func Binary(key string, val []byte) Field {
	return zap.Binary(key, val)
}

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

func ByteString(key string, val []byte) Field {
	return zap.ByteString(key, val)
}

func Complex128(key string, val complex128) Field {
	return zap.Complex128(key, val)
}

func Complex64(key string, val complex64) Field {
	return zap.Complex64(key, val)
}

func String(key string, val string) Field {
	return zap.String(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Int8(key string, val int8) Field {
	return zap.Int8(key, val)
}

func Int16(key string, val int16) Field {
	return zap.Int16(key, val)
}

func Int32(key string, val int32) Field {
	return zap.Int32(key, val)
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Uint(key string, val uint) Field {
	return zap.Uint(key, val)
}

func Uint8(key string, val uint8) Field {
	return zap.Uint8(key, val)
}

func Uint16(key string, val uint16) Field {
	return zap.Uint16(key, val)
}

func Uint32(key string, val uint32) Field {
	return zap.Uint32(key, val)
}

func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}

func Uintptr(key string, val uintptr) Field {
	return zap.Uintptr(key, val)
}

func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func Reflect(key string, val interface{}) Field {
	return zap.Reflect(key, val)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

// Stack constructs a field that stores a stacktrace of the current goroutine
// under provided key. Keep in mind that taking a stacktrace is eager and
// expensive (relatively speaking); this function both makes an allocation and
// takes about two microseconds.
func Stack(key string) Field {
	return zap.Stack(key)
}

func NamedError(err error) Field {
	return zap.Error(err)
}

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
//
// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
// them. To minimize surprises, []byte values are treated as binary blobs, byte
// values are treated as uint8, and runes are always treated as integers.
func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}

// Namespace creates a named, isolated scope within the logger's context. All
// subsequent fields will be added to the new namespace.
//
// This helps prevent key collisions when injecting loggers into sub-components
// or third-party libraries.
func Namespace(key string) Field {
	return zap.Namespace(key)
}

// Stringer constructs a field with the given key and the output of the value's
// String method. The Stringer's String method is called lazily.
func Stringer(key string, val fmt.Stringer) Field {
	return zap.Stringer(key, val)
}

func Sync() {
	zapLogger.Sync()
}

func Debug(msg string, fields ...Field) {
	zapLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	zapLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	zapLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	zapLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	zapLogger.Fatal(msg, fields...)
}
