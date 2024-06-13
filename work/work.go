package work

import (
	"errors"
	"fmt"
	"github.com/orestonce/m3u8d"
	"github.com/orestonce/m3u8d/m3u8dcpp"
	"github.com/sirupsen/logrus"
	"gom3u8/data"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type Work struct {
	ID         int       `json:"ID"`
	Name       string    `json:"name"`
	Url        string    `json:"url"`
	SaveDir    string    `json:"save_dir"`
	State      int       `json:"state"`
	Info       string    `json:"info"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (w Work) TableName() string {
	return "work_info"
}

func (work *Work) Save(url string, fileName string, save_dir string) (err error) {
	url = strings.Replace(url, " ", "", -1)
	fileName = strings.Replace(fileName, " ", "", -1)
	save_dir = strings.Replace(save_dir, " ", "", -1)

	if fileName == "" {
		fileName, err = extractFilenameFromURL(url)
		if err != nil {
			return err
		}
	}
	if save_dir == "" {
		save_dir = "/data/m3u8/"
	}
	fmt.Println("url:", url)
	fmt.Println("fileName:", fileName)
	fmt.Println("save_dir:", save_dir)
	db := data.DataDB
	workInfo := &Work{
		Name:       fileName,
		Url:        url,
		State:      1,
		SaveDir:    save_dir,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	db = db.Model(workInfo)
	db = db.Where("url=?", url).FirstOrCreate(&workInfo)

	if db.Error != nil {
		fmt.Println("save err", db.Error)
		return db.Error
	}
	fmt.Println("save:", workInfo)
	return
}

func (work *Work) GetNotWorkingWork() *Work {
	workInfo := &Work{
		State: 1,
	}
	db := data.DataDB
	db.Error = nil
	db.Model(work)
	db = db.First(workInfo, "state = 1")

	if db.Error != nil {
		logrus.Error(db.Error)
		return nil
	}
	return workInfo
}
func (work *Work) End() error {
	db := data.DataDB
	db.Model(work)
	work.State = 2
	db.Save(work)
	return nil
}
func (work *Work) Error(err_msg string) error {
	db := data.DataDB
	db.Model(work)
	work.State = 3
	work.Info = err_msg
	db.Save(work)
	return nil
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

func DownloadFromCmd(req m3u8d.StartDownload_Req) error {
	req.ProgressBarShow = true
	fmt.Println(req.M3u8Url)
	errMsg := m3u8dcpp.StartDownload(req)
	if errMsg != "" {
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}

	resp := m3u8dcpp.WaitDownloadFinish()

	if resp.ErrMsg != "" {
		fmt.Println(resp.ErrMsg)
		return errors.New(errMsg)
	}
	if resp.IsSkipped {
		fmt.Println("已经下载过了: " + resp.SaveFileTo)
		return errors.New(errMsg)
	}
	if resp.SaveFileTo == "" {
		fmt.Println("下载成功.")
		return errors.New(errMsg)
	}
	fmt.Println("下载成功, 保存路径", resp.SaveFileTo)
	return nil
}
