package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gom3u8/util"
	"math"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next() // 调用该请求的剩余处理程序
		stopTime := time.Since(startTime)
		spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds()/1000000))))
		//hostName, err := os.Hostname()
		//if err != nil {
		//    hostName = "Unknown"
		//}
		statusCode := c.Writer.Status()
		//clientIP := c.ClientIP()
		//userAgent := c.Request.UserAgent()
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method
		url := c.Request.RequestURI
		Log := util.Logger.WithFields(logrus.Fields{
			//"HostName": hostName,
			"SpendTime": spendTime,
			"path":      url,
			"Method":    method,
			"status":    statusCode,
			//"Ip": clientIP,
			//"DataSize": dataSize,
			//"UserAgent": userAgent,
		})
		if len(c.Errors) > 0 { // 矿建内部错误
			Log.Error(c.Errors.ByType(gin.ErrorTypePrivate))
		}
		if statusCode >= 500 {
			Log.Error()
		} else if statusCode >= 400 {
			Log.Warn()
		} else {
			Log.Info()
		}
	}
}
