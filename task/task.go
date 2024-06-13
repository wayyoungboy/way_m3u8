package task

import (
	"github.com/gin-gonic/gin"
	"gom3u8/work"
)

type TaskController struct {
}

func (TaskController) AddTask(c *gin.Context) {
	url := c.PostForm("url")
	save_dir := c.PostForm("save_dir")
	file_name := c.PostForm("file_name")
	if len(save_dir) == 0 || len(file_name) == 0 {
		save_dir = "/data1/media/av"
	}
	w := &work.Work{}
	err := w.Save(url, file_name, save_dir)
	if err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "ok"})
}
