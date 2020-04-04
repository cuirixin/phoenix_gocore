## 说明
日志


## Usage

```go
// 自定义Gin日志
router.Use(commonlog.LoggerWithFormatter())
// 自定义日志，按日期切割
router.Use(commonlog.LoggerToFileWithRes())
// 自定义日志，按日期切割，带请求参数和输出结果
router.Use(commonlog.LoggerToFileWithReqRes(conf.Conf.Log.ApiFilePath, conf.Conf.Log.ApiFileName))
```