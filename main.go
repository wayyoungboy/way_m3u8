package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gom3u8/conf"
	_ "gom3u8/data"
	"gom3u8/task"
	"gom3u8/work"
	"gopkg.in/natefinch/lumberjack.v2"
	"strconv"
)

func main() {

	conf.ConfInit()
	log_nu, err := strconv.Atoi(conf.ConfMap["log_Nu"].(string))
	if err != nil {
		fmt.Println("log_Nu err:", err)
		return
	}
	logFile := &lumberjack.Logger{
		Filename:   "./log/log.txt",
		MaxSize:    10, // MB
		MaxBackups: log_nu,
		MaxAge:     28, // days
		Compress:   true,
		LocalTime:  true,
	}

	defer logFile.Close()

	log.SetOutput(logFile)

	run()

}
func run() {
	go work.Working()
	r := gin.Default()
	tc := task.TaskController{}
	r.Static("/static", "./static")
	r.POST("/addTask", tc.AddTask)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})
	fmt.Println("open http://127.0.0.1:2045/static/ ")
	r.Run(":" + fmt.Sprint(conf.ConfMap["Init.Port"])) // 监听2045端口
}
