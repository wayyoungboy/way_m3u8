package data

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DataDB *gorm.DB

func init() {
	var err error
	// 调用 Open 方法，传入驱动名和连接字符串
	DataDB, err = gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	// 检查是否有错误
	if err != nil {
		fmt.Println("连接数据库失败：", err)
		return
	}
	// 打印成功信息
	fmt.Println("连接数据库成功")
}
