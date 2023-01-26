package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yhmain/5th-simple-tiktok/config"
)

func main() {
	r := gin.Default()
	initRouter(r)                      // 初始化路由
	params := config.ProjectConfig.App // 获取配置文件中的 主机IP和端口号
	addr := fmt.Sprintf("%s:%s", params.Host, params.Port)
	r.Run(addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
