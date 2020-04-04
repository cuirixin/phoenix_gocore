package commonlog

import (
	"bytes"
	putils "github.com/cuirixin/phoenix_gocore/utils"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

var logger *logrus.Logger

func Info (module, msg string)  {
	logger.WithFields(logrus.Fields{
		"module"  : module,
		"msg"     : msg,
	}).Info()
}

func initLogger(logFilePath, logFileName string) {
	//logFilePath := conf.Conf.Log.ApiFilePath
	//logFileName := conf.Conf.Log.ApiFileName

	fileName := path.Join(putils.CallerSourcePath(), logFilePath, logFileName)
	color.Green("日志文件路径: %s", fileName)

	// 禁止logrus的输出
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		color.Red("打开日志文件失败，#{err}")
	}

	// 实例化
	logger = logrus.New()
	// 设置输出
	logger.Out = src
	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName + ".%Y-%m-%d-%H-%M.log",
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.TextFormatter{ // JsonFormatter
		TimestampFormat:"2006-01-02 15:04:05",
	})

	// 新增 Hook
	logger.AddHook(lfHook)

}

// 日志记录到文件
func LoggerToFile(logFilePath, logFileName string) gin.HandlerFunc {

	initLogger(logFilePath, logFileName)

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		entry := logger.WithFields(logrus.Fields{
			"status"  : statusCode,
			"latency" : latencyTime,
			"ip"      : clientIP,
			"method"  : reqMethod,
			"uri"     : reqUri,
		})
		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			entry.Error(c.Errors.String())
		} else {
			entry.Info()
		}
	}
}


const MAX_PRINT_BODY_LEN = 512

type bodyLogWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	//memory copy here!
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}

// 日志记录到文件
func LoggerToFileWithReqRes(logFilePath, logFileName string) gin.HandlerFunc {

	initLogger(logFilePath, logFileName)

	return func(c *gin.Context) {
		req_body, _ := ioutil.ReadAll(c.Request.Body)

		// 开始时间
		startTime := time.Now()
		// 初始化bodyLogWriter
		strBody := ""
		var blw bodyLogWriter
		//if we need to log res body
		blw = bodyLogWriter{bodyBuf: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 返回消息体
		strBody = strings.Trim(blw.bodyBuf.String(), "\n")
		if len(strBody) > MAX_PRINT_BODY_LEN {
			strBody = strBody[:(MAX_PRINT_BODY_LEN - 1)]
		}
		// 日志格式
		entry := logger.WithFields(logrus.Fields{
			"status"  : statusCode,
			"latency" : latencyTime,
			"ip"      : clientIP,
			"method"  : reqMethod,
			"uri"     : reqUri,
			"req_body"     : string(req_body),
			"res"     : strBody,
		})
		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			entry.Error(c.Errors.String())
		} else {
			entry.Info()
		}
	}
}

// 日志记录到 MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
