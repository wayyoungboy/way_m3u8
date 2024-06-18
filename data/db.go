package data

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var DataDB *gorm.DB

func init() {
	var err error
	slowLogger := logger.New(
		//设置Logger
		NewMyWriter(),
		logger.Config{
			//慢SQL阈值
			SlowThreshold: time.Millisecond,
			//设置日志级别，只有Warn以上才会打印sql
			LogLevel: logger.Error,
		},
	)
	// 调用 Open 方法，传入驱动名和连接字符串
	DataDB, err = gorm.Open(sqlite.Open("./workinfo.db"), &gorm.Config{
		Logger: slowLogger,
	})
	// 检查是否有错误
	if err != nil {
		fmt.Println("连接数据库失败：", err)
		return
	}
	// 打印成功信息
	fmt.Println("连接数据库成功")

}

type gormLog struct {
}

func NewGormLogger() *gormLog {
	return &gormLog{}
}

// 定义自己的Writer
type MyWriter struct {
	mlog *logrus.Logger
}

// 实现gorm/logger.Writer接口
func (m *MyWriter) Printf(format string, v ...interface{}) {
	logstr := fmt.Sprintf(format, v...)
	//利用loggus记录日志
	m.mlog.Info(logstr)
}

func NewMyWriter() *MyWriter {
	log := logrus.New()
	//配置logrus
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return &MyWriter{mlog: log}
}
