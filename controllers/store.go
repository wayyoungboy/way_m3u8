package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orestonce/m3u8d"
	"net/url"
	"path/filepath"
	"way_m3u8/worker"
)

func StoreURL(c *gin.Context) {
	fmt.Println("DD")

	url := c.PostForm("url")
	save_dir := c.PostForm("save_dir")
	file_name := c.PostForm("file_name")
	if len(save_dir) == 0 || len(file_name) == 0 {
		save_dir = "/data1/media/av"
	}
	if len(file_name) == 0 || len(file_name) == 0 {
		fileName, err := extractFilenameFromURL(url)
		if err != nil {
			c.JSON(400, gin.H{"err": err.Error()})
			return
		}
		if fileName == "" {
			err := errors.New("can not find file name by url")
			c.JSON(400, gin.H{"err": err.Error()})
			return
		} else {
			file_name = fileName
		}
	}
	fmt.Println("filename:", file_name)
	fmt.Println("save_dir:", save_dir)

	req := m3u8d.StartDownload_Req{
		M3u8Url:                  url,
		Insecure:                 false,
		SaveDir:                  save_dir,
		FileName:                 file_name,
		SkipTsExpr:               "",
		SetProxy:                 "",
		HeaderMap:                nil,
		SkipRemoveTs:             false,
		ProgressBarShow:          false,
		ThreadCount:              8,
		SkipCacheCheck:           false,
		SkipMergeTs:              false,
		Skip_EXT_X_DISCONTINUITY: false,
		DebugLog:                 false,
	}
	fmt.Println(req)
	if err := worker.AddWork(req); err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "URL stored successfully"})
}
func extractFilenameFromURL(urlStr string) (string, error) {
	// 解析URL
	_, err := url.ParseRequestURI(urlStr)
	// 如果没有错误，则认为URL是合法的
	if err != nil {
		return "", err
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	// 使用path/filepath包来获取路径的最后一个元素，这将是文件名
	filename := filepath.Base(parsedURL.Path)
	if len(filename) > 5 {
		filename = filename[:len(filename)-5]
	}
	return filename, nil
}
