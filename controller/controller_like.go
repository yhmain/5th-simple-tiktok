package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/yhmain/5th-simple-tiktok/dao"
	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/util"
)

func init() {
	JobSaveRedis() // 开启定时器任务
}

/*
	关于赞的操作有点赞和取消点赞
	数据流：
	1. 修改redis中的数据为1: true   0: false
	2. (定时任务)保存或者是更新到mysql
*/

//用户赞操作函数，路由
func FavoriteAction(c *gin.Context) {
	// 经过jwt中间件
	usertoken := c.MustGet("usertoken").(middleware.UserToken)
	video_id := middleware.GetParamPostOrGet(c, "video_id")
	action_type := c.Query("action_type")
	redisKey := fmt.Sprintf("Fav:%d:%s", usertoken.UserID, video_id)      // 加上前缀Fav:，取用户id，作为redis的key
	_, err := middleware.UpdateRedisLike(redisKey, video_id, action_type) // 更新Redis
	if err != nil {                                                       // 赞操作失败
		fmt.Println(err)
		c.JSON(http.StatusOK, util.FavActionErr)
		return
	}
	// fmt.Println("点赞成功!", action_type, " ** ", video_id, "Redis key: ", redisKey)
	// 保存到MySQL
	// 实际上利用Redis+定时任务即可
	// SaveRedisToMySQL() // 此处是直接保存，另外一种是设置定时任务
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
	likes, err := middleware.ParseRedisKeys()
	if err != nil {
		fmt.Println("解析Redis的Keys出现异常：", err)
		return
	}
	// 调用dao层保存到数据库
	if len(likes) < 1 {
		// fmt.Println("空点赞数据，无须保存.")
		return
	}
	// 赞数据插入到数据库
	err = dao.UpsetLikes(&likes)
	if err != nil {
		fmt.Println("批量更新点赞数据出错：", err)
		return
	}
	// 之后清空Redis缓存
	middleware.ClearRedis()
}

//用户喜欢列表函数，路由
func FavoriteList(c *gin.Context) {
	// 经过jwt中间件
	usertoken := c.MustGet("usertoken").(middleware.UserToken)
	// 查询之前，将Redis数据手段更新到Mysql
	SaveRedisToMySQL()
	// 此后，直接查询Redis数据
	likes := dao.GetUserLike(usertoken.UserID) // 获取到了视频ID
	vids := make([]int64, len(likes))          // 提取出来存在切片里面
	for i := range likes {
		vids[i] = likes[i].VideoID
	}
	// 调用Dao层 进行批量查询
	videos := dao.GetVideoByIDs(vids)
	// 返回
	c.JSON(http.StatusOK, FeedResponse{
		Response:  util.Success, //成功
		VideoList: videos,
	})
}
