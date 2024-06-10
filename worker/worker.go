package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/orestonce/m3u8d"
	"github.com/orestonce/m3u8d/m3u8dcpp"
	"sync"
	"time"
)

var Work_list_max = 20
var Work_list *MessageQueue

func AddWork(workinfo m3u8d.StartDownload_Req) error {
	if Work_list.Len() >= Work_list_max {
		return errors.New("Work_list.Len()>=Work_list_max")
	}
	Work_list.Push(workinfo)
	return nil
}
func WorkInit() {
	Work_list = NewMessageQueue()

}
func WorkRun() {
	for {
		if Work_list.Len() > 0 {
			work := Work_list.Pop()
			downloadFromCmd(work)
		} else {
			time.Sleep(5 * time.Second)
		}

	}
}
func downloadFromCmd(req m3u8d.StartDownload_Req) {
	req.ProgressBarShow = true
	errMsg := m3u8dcpp.StartDownload(req)
	if errMsg != "" {
		fmt.Println(errMsg)
		return
	}

	e, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(e))
	resp := m3u8dcpp.WaitDownloadFinish()
	fmt.Println() // 有进度条,所以需要换行
	fmt.Println("dd")
	if resp.ErrMsg != "" {
		fmt.Println(resp.ErrMsg)
		return
	}
	if resp.IsSkipped {
		fmt.Println("已经下载过了: " + resp.SaveFileTo)
		return
	}
	if resp.SaveFileTo == "" {
		fmt.Println("下载成功.")
		return
	}
	fmt.Println("下载成功, 保存路径", resp.SaveFileTo)
}

// MessageQueue 是一个简单的消息队列，使用通道实现
type MessageQueue struct {
	queue chan m3u8d.StartDownload_Req
	mu    sync.Mutex // 可选，如果需要并发安全
}

// NewMessageQueue 创建一个新的消息队列
func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		queue: make(chan m3u8d.StartDownload_Req, Work_list_max), // 容量为Work_list_max的缓冲通道
	}
}

// Push 向队列中添加消息
func (q *MessageQueue) Push(message m3u8d.StartDownload_Req) {
	// 由于channel本身是并发安全的，这里的锁可以省略
	// 但如果你需要保护其他共享状态，可以加上锁
	q.queue <- message
}

// Pop 从队列中取出消息（如果队列为空，会阻塞）
func (q *MessageQueue) Pop() m3u8d.StartDownload_Req {
	return <-q.queue
}

// Len 返回队列中消息的数量（注意：这不是并发安全的）
func (q *MessageQueue) Len() int {
	return len(q.queue)
}

type WorkInfo struct {
	Id       int32  `json:"id"`
	Url      string `json:"url"`
	FileName string `json:"file_name"`
	State    string `json:"state"`
}
