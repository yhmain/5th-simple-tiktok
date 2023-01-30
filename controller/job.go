package controller

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/yhmain/5th-simple-tiktok/dao"
	"github.com/yhmain/5th-simple-tiktok/middleware"
)

func init() {
	JobSaveRedis() // 开启定时器任务
}

// 定时任务：将Redis数据保存到MySQL
func JobSaveRedis() {
	// 启动点赞操作的定时任务
	// 定时将Redis的数据保存到数据库
	c := cron.New()                  // 这个对象用于管理定时任务
	c.AddFunc("@every 10s", func() { // @every后加一个时间间隔，表示每隔多长时间触发一次 eg. 2s  1m2s 1h
		fmt.Println(time.Now(), "Tick every 10 second: save redis data to mysql.")
		SaveRedisToMySQL()
	})
	c.Start() // 启动定时循环
}

// 将Redis数据保存到MySQL
func SaveRedisToMySQL() {
	// 赞信息和赞数量，更新到数据库
	likes, likeCount := middleware.GetRedisLike()           // 调用redis解析里面的keys
	if err := dao.SaveLikes(likes, likeCount); err != nil { // 赞数据插入到数据库
		fmt.Println("批量更新点赞数据出错：", err)
	}
	// 接下来是 评论内容和评论数量，更新到数据库
	commentAdd, commendDel, commentCnt := middleware.GetRedisComment()           // 调用redis解析里面的keys
	if err := dao.SaveComments(commentAdd, commendDel, commentCnt); err != nil { // 评论数据插入到数据库
		fmt.Println("批量更新评论数据出错：", err)
	}
	// 接下来是 关注操作的信息，更新到数据库
	follows, followCount := middleware.GetRedisFollow()
	if err := dao.SaveFollows(follows, followCount); err != nil { // 评论数据插入到数据库
		fmt.Println("批量更新关注数据出错：", err)
	}
}
