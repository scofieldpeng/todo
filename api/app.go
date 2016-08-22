package main

import (
	"flag"
	"github.com/labstack/echo/engine/standard"
	"github.com/scofieldpeng/config-go"
	"github.com/scofieldpeng/mysql-go"
	"github.com/scofieldpeng/redis-go"
	"github.com/scofieldpeng/todo/api/libs/common"
	_ "github.com/scofieldpeng/todo/api/libs/log"
	_ "github.com/scofieldpeng/todo/api/routes"
	"log"
	"github.com/scofieldpeng/todo/api/libs/auth"
	"github.com/scofieldpeng/template-go"
	"path/filepath"
	"os"
	"github.com/tylerb/graceful"
	"time"
)

func init() {
	flag.BoolVar(&common.Debug, "debug", false, "是否是debug模式,默认false")
	flag.Parse()

	common.Echo.SetDebug(common.Debug)

	if err := config.New("", common.Debug); err != nil {
		log.Fatalln("初始化config失败!错误原因:", err)
	}
	log.Println("初始化config完成")

	// 初始化auth
	if err := auth.Init();err != nil {
		log.Fatalln("初始化auth失败,错误原因:",err)
	}
	log.Println("初始化auth成功")

	// 初始化redis
	redis.Init(config.Config("redis"))
	log.Println("初始化redis完成")

	// 初始化myql
	if err := mysql.Init(config.Config("mysql"),common.Debug); err != nil {
		log.Fatalln("初始化mysql失败,错误原因:", err)
	}
	log.Println("初始化mysql完成")

	// 初始化tpl
	template.Tpl.SetTplSuffix(".html")
	template.Tpl.SetDelimeter("[[","]]")
	currentPath,_ := filepath.Abs(filepath.Dir(os.Args[0]))
	if err := template.Tpl.New(currentPath + string(os.PathSeparator) + "tpls" + string(os.PathSeparator));err != nil {
		log.Fatalln("初始化tpl失败!错误原因:",err.Error())
	}
	common.Echo.SetRenderer(template.Tpl)
}

func main() {
	host,_ := config.Config("app").Get("app","host")
	port,_ := config.Config("app").Get("app","port")
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "8081"
	}

	std := standard.New(host + ":" + port)
	std.SetHandler(common.Echo)
	graceful.ListenAndServe(std.Server,5*time.Second)
}
