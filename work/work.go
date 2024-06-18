package work

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/orestonce/m3u8d"
	"github.com/orestonce/m3u8d/m3u8dcpp"
	log "github.com/sirupsen/logrus"
	"gom3u8/conf"
	"gom3u8/data"
	"gorm.io/gorm"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	StateAll   = 0
	StateReady = 1
	StateEnd   = 2
	StateError = 3
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
	err = db.Where("url=?", url).First(&workInfo).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

	}
	if workInfo.ID != 0 {
		id := fmt.Sprint(workInfo.ID)
		log.Warn("work already exists: " + id)
		return errors.New("work already exists: " + id)
	}

	db = db.Model(workInfo)
	err = db.Where("url=?", url).Create(&workInfo).Error
	if err != nil {

		return err
	}

	log.Info("save:", workInfo)
	return
}

func (work *Work) GetNotWorkingWork(startId int) *Work {

	workInfoList := []Work{}
	db := data.DataDB
	db.Error = nil
	db.Model(work)
	db = db.Where("state = ?", StateReady).Order("id").Where("id>?", startId).Find(&workInfoList)
	if db.RowsAffected > 0 {
		return &workInfoList[0]
	}
	return nil
}
func (work *Work) List(Limit, Offset int) *[]Work {
	db := data.DataDB
	db.Error = nil
	db.Model(work)
	return nil
}
func (work *Work) End() error {
	db := data.DataDB
	db.Model(work)
	work.State = StateEnd
	db.Save(work)
	return nil
}
func (work *Work) Error(err_msg string) error {
	log.Warn(work.ID, " download err:", err_msg)
	db := data.DataDB
	db.Model(work)
	work.State = StateError
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
	workListMaxNu := conf.ConfMap["work_max"].(int)
	workList := []Worker{}
	for i := 0; i < workListMaxNu; i++ {
		workList = append(workList, NewWorker())
	}
	for {
		w := &Work{}
		//过滤已处理的任务
		startId := 0
		for {
			for _, worker := range workList {
				if worker.State {
					readywork := w.GetNotWorkingWork(startId)
					if readywork == nil {
						time.Sleep(5 * time.Second)
						continue
					}
					startId = readywork.ID
					worker.State = false
					go worker.Start(readywork)
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}

type Worker struct {
	State bool
}

func NewWorker() Worker {
	return Worker{
		State: true,
	}
}
func (w Worker) Start(readywork *Work) {
	workOnece(readywork)
	w.State = true
}

func workOnece(readywork *Work) {
	log.Info("readywork:", readywork)
	_, err := os.Stat(readywork.SaveDir)
	if err != nil {
		os.MkdirAll(readywork.SaveDir, os.ModePerm)
	}
	domain, err := ExtractDomain(readywork.Url)

	if err != nil {
		log.Warn("域名解析错误： ", err)
		return
	}

	header_filename := "header.json"
	headers_list := []URLData{}

	header_data, err := os.ReadFile(header_filename)
	if err != nil {
		readywork.Error("Error reading header jspn file:" + err.Error())
		return
	}
	err = json.Unmarshal(header_data, &headers_list)
	if err != nil {
		log.Error("Error parsing JSON:", err)
		readywork.Error("Error parsing JSON:" + err.Error())
		return
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
		return
	}

	readywork.End()

}

// 定义header.json主结构体
type URLData struct {
	URL    string              `json:"url"`
	Header map[string][]string `json:"header"`
}
