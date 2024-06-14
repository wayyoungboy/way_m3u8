package work

import (
	"encoding/json"
	"errors"
	"github.com/orestonce/m3u8d"
	"github.com/orestonce/m3u8d/m3u8dcpp"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gom3u8/data"
	"net/url"
	"os"
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
	log.Info("url:", url)
	log.Info("fileName:", fileName)
	log.Info("save_dir:", save_dir)
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
		log.Error("save err", db.Error)
		return db.Error
	}
	log.Info("save:", workInfo)
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
	log.Warn(work.ID, " download err:", err_msg)
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
	log.Info("DownloadFromCmd M3u8Url: ", req.M3u8Url)
	errMsg := m3u8dcpp.StartDownload(req)
	if errMsg != "" {
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	resp := m3u8dcpp.WaitDownloadFinish()

	if resp.ErrMsg != "" {
		log.Error(resp.ErrMsg)
		return errors.New(errMsg)
	}
	if resp.IsSkipped {
		log.Warn("已经下载过了: " + resp.SaveFileTo)
		return errors.New(errMsg)
	}
	if resp.SaveFileTo == "" {
		log.Info("下载成功.")
		return errors.New(errMsg)
	}
	log.Info("下载成功, 保存路径", resp.SaveFileTo)
	return nil
}
func ExtractDomain(rawURL string) (string, error) {
	// 解析URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// 获取主机名
	hostname := parsedURL.Hostname()

	// 如果主机名包含端口号，则去除它
	if strings.Contains(hostname, ":") {
		parts := strings.SplitN(hostname, ":", 2)
		hostname = parts[0]
	}

	return hostname, nil
}

func Working() {
	for {
		w := &Work{}
		readywork := w.GetNotWorkingWork()
		log.Info("readywork:", readywork)
		//没有任务的时候停五秒再重复
		if readywork == nil {
			time.Sleep(5 * time.Second)
			continue
		}
		_, err := os.Stat(readywork.SaveDir)
		if err != nil {
			os.MkdirAll(readywork.SaveDir, os.ModePerm)
		}
		domain, err := ExtractDomain(readywork.Url)

		if err != nil {
			log.Warn("域名解析错误： ", err)
			continue
		}

		header_filename := "header.json"
		headers_list := []URLData{}

		header_data, err := os.ReadFile(header_filename)
		if err != nil {
			log.Error("Error reading header jspn file:", err)
			return
		}
		err = json.Unmarshal(header_data, &headers_list)
		if err != nil {
			log.Error("Error parsing JSON:", err)
			continue
		}
		header := make(map[string][]string)
		for _, headers_data := range headers_list {
			if headers_data.URL == domain {
				header = headers_data.Header
			}
		}
		req := m3u8d.StartDownload_Req{
			M3u8Url:                  readywork.Url,
			Insecure:                 true,
			SaveDir:                  readywork.SaveDir,
			FileName:                 readywork.Name,
			SkipTsExpr:               "",
			SetProxy:                 "",
			HeaderMap:                header,
			SkipRemoveTs:             false,
			ProgressBarShow:          false,
			ThreadCount:              8,
			SkipCacheCheck:           false,
			SkipMergeTs:              false,
			Skip_EXT_X_DISCONTINUITY: false,
			DebugLog:                 false,
		}
		err = DownloadFromCmd(req)

		if err != nil {

			readywork.Error(err.Error())
			continue
		}

		readywork.End()
	}

}

// 定义header.json主结构体
type URLData struct {
	URL    string              `json:"url"`
	Header map[string][]string `json:"header"`
}
