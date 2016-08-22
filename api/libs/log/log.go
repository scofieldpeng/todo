package log

import (
	"fmt"
	"github.com/scofieldpeng/todo/api/libs/common"
	"log"
	"os"
	"time"
)

// init 初始化log包
func init() {
	log.SetPrefix("[todo]")
	log.SetFlags(log.Ltime | log.Llongfile)

	go func() {
		for {
			update()
		}
	}()
}

// update 更新log输出的文件名称
func update() {
	nowTime := time.Now()
	year, month, day := nowTime.Date()
	nowTimestamp := nowTime.Unix()
	tomorrowTimestamp := time.Date(year, month, day+1, 0, 0, 0, 0, nowTime.Location()).Unix()
	setOutput(fmt.Sprintf("%d-%d-%d.log", year, int(month), day))

	time.Sleep(time.Duration(tomorrowTimestamp-nowTimestamp) * time.Second)
}

// setOutput 更新log输出的文件文件
func setOutput(fileName string) {
	// 检查log文件夹是否存在
	os.Mkdir(common.AppDir()+string(os.PathSeparator)+"log"+string(os.PathSeparator), os.FileMode(0755))
	obj, err := os.OpenFile(common.AppDir()+string(os.PathSeparator)+"log"+string(os.PathSeparator)+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0655))
	if err != nil {
		log.Println("设置log的out文件出错,原因:", err.Error())
		log.SetOutput(os.Stdout)
	}
	log.SetOutput(obj)
}
