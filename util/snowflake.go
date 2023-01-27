package util

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
	"github.com/yhmain/5th-simple-tiktok/config"
)

//分布式id生成器
var userIDGen *snowflake.Node

// 先进行初始化
func init() {
	fmt.Println("snowflake.go ... init")
	params := config.GetConfig().Snowflake
	node, err := snowflake.NewNode(params.MachineID)
	if err != nil {
		fmt.Println("userIDGen: ", node)
		fmt.Println("err: ", err)
		return
	}
	userIDGen = node // 赋值给全局变量
	fmt.Println("snowflake.go ... 初始化成功！")
}

// 生成用户ID
func GenID() int64 {
	fmt.Println("userIDGen是否为空: ", userIDGen == nil)
	return userIDGen.Generate().Int64()
}
