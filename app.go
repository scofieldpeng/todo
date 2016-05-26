package main

import (
	"flag"
	"github.com/labstack/echo/engine/standard"
	"github.com/scofieldpeng/config-go"
	"github.com/scofieldpeng/mysql-go"
	"github.com/scofieldpeng/redis-go"
	"github.com/scofieldpeng/todo/libs/common"
	_ "github.com/scofieldpeng/todo/libs/log"
	_ "github.com/scofieldpeng/todo/routes"
	"log"
	"github.com/scofieldpeng/todo/libs/auth"
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
	if err := mysql.Init(config.Config("mysql")); err != nil {
		log.Fatalln("初始化mysql失败,错误原因:", err)
	}
	log.Println("初始化mysql完成")

}

func main() {
	common.Echo.Run(standard.New(":8081"))
}
