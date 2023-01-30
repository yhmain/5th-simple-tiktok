package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yhmain/5th-simple-tiktok/dao"
	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/util"
)

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
	redisKey := fmt.Sprintf("Fav:%d:%s", usertoken.UserID, video_id) // 加上前缀Fav:，取用户id，作为redis的key
	val, err := middleware.GetKey(redisKey)                          // 先获取当前 点赞的状态
	// 什么情况不用更改状态？ Key命中且当前操作类型与value代表含义一样
	if (err == nil && action_type == "1" && val == "1") || (err == nil && action_type == "2" && val == "0") {
		c.JSON(http.StatusOK, util.FavActionErr)
		return
	}
	redisLikeCnt := fmt.Sprintf("FavCnt:%s", video_id) // 视频赞数量的Key
	if !middleware.ExistKey(redisLikeCnt) {
		vID, _ := strconv.ParseInt(video_id, 10, 64)
		video := dao.GetVideoByID(vID)                                              // 此处是查询数据库
		middleware.SetKey(redisLikeCnt, strconv.FormatInt(video.FavoriteCount, 10)) // 设置key
	}
	err = middleware.UpdateRedisLike(redisKey, video_id, action_type) // 更新Redis
	if err != nil {                                                   // 赞操作失败
		fmt.Println(err)
		c.JSON(http.StatusOK, util.FavActionErr)
		return
	}
	// fmt.Println("点赞成功!", action_type, " ** ", video_id, "Redis key: ", redisKey)
	// 保存到MySQL
	// 实际上利用Redis+定时任务即可
	// SaveRedisToMySQL() // 此处是直接保存，另外一种是设置定时任务
}

//用户喜欢列表函数，路由
func FavoriteList(c *gin.Context) {
	// 经过jwt中间件
	usertoken := c.MustGet("usertoken").(middleware.UserToken)
	// 查询之前，将Redis数据手动更新到Mysql
	SaveRedisToMySQL()
	// 再查询MySQL数据
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
