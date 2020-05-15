package commonlog

import (
	"bytes"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"mime"
	"path"
	"strings"
	"time"
)

var logger *logrus.Logger

func Log() *logrus.Logger {
	if logger == nil {
		initDefaultLogger()
	}
	return logger
}

func initDefaultLogger()  {
	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.AddHook(NewContextHook())
}

func initLogger(logFilePath, logFileName string) {
	//logFilePath := conf.Conf.Log.ApiFilePath
	//logFileName := conf.Conf.Log.ApiFileName

	fileName := path.Join(logFilePath, logFileName)
	color.Green("日志文件路径: %s", fileName)

	// 实例化
	logger = logrus.New()
	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	// 新增行号和文件名Hook
	logger.AddHook(NewContextHook())
	// 自动切割文件Hook
	logger.AddHook(NewLfsHook(fileName, 1000))



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
			"req_id"  : c.Writer.Header().Get("X-Request-Id"), // 依赖于request_id中间件
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

		// 打印RequestBody，目前只针对JSON格式数据

		var reqBody []byte
		ctGet := c.Request.Header.Get("Content-Type")
		ct, _, _ := mime.ParseMediaType(ctGet)
		switch ct {
		case gin.MIMEJSON:
			buf, _ := ioutil.ReadAll(c.Request.Body)
			rdr := ioutil.NopCloser(bytes.NewBuffer(buf))
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
			reqBody, _ = ioutil.ReadAll(rdr)
		}

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
			"req_id"  : c.Writer.Header().Get("X-Request-Id"), // 依赖于request_id中间件
			"req_body": string(reqBody),
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
