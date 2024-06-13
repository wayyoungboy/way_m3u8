package main

import (
	"fmt"
	"github.com/orestonce/m3u8d"
	_ "gom3u8/data"
	"gom3u8/middleware"
	"gom3u8/task"
	"gom3u8/work"
	"io"
	"os"
	"time"
)
import "github.com/gin-gonic/gin"

func main() {
	run()

}
func run() {
	logName := "./log/gin.log"
	_, err := os.Stat(logName)
	var f *os.File
	if err != nil {
		f, _ = os.Create(logName)
	}
	f, err = os.Open(logName)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	go working()
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())
	tc := task.TaskController{}
	r.Static("/static", "./static")
	r.POST("/addTask", tc.AddTask)
	r.Run(":2045") // 监听2045端口

}
func working() {
	for {
		w := &work.Work{}
		readywork := w.GetNotWorkingWork()
		fmt.Println("readywork:", readywork)
		if readywork == nil {
			time.Sleep(5 * time.Second)
			continue
		}
		_, err := os.Stat(readywork.SaveDir)
		if err != nil {
			os.MkdirAll(readywork.SaveDir, os.ModePerm)
		}
		m := make(map[string][]string)
		m["Referer"] = []string{}
		req := m3u8d.StartDownload_Req{
			M3u8Url:                  readywork.Url,
			Insecure:                 true,
			SaveDir:                  readywork.SaveDir,
			FileName:                 readywork.Name,
			SkipTsExpr:               "",
			SetProxy:                 "",
			HeaderMap:                m,
			SkipRemoveTs:             false,
			ProgressBarShow:          false,
			ThreadCount:              8,
			SkipCacheCheck:           false,
			SkipMergeTs:              false,
			Skip_EXT_X_DISCONTINUITY: false,
			DebugLog:                 false,
		}
		err = work.DownloadFromCmd(req)
		fmt.Println(err)
		if err != nil {
			readywork.Error(err.Error())
			continue
		}

		readywork.End()
	}

}
