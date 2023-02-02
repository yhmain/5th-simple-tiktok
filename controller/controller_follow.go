package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yhmain/5th-simple-tiktok/dao"
	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
	"github.com/yhmain/5th-simple-tiktok/util"
)

// 关注/粉丝功能 的响应结构体
type FollowListResponse struct {
	util.Response
	UserList []model.User `json:"user_list"`
}

//关注操作函数，路由
func RelationAction(c *gin.Context) {
	usertoken := c.MustGet("usertoken").(middleware.UserToken) // 经过jwt中间件
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	ubID, _ := strconv.ParseInt(to_user_id, 10, 64)
	// 首先防止重复操作，什么情况不用更改状态？ Key命中且当前操作类型与value代表含义一样
	followKey := fmt.Sprintf("Fol:%d:%d", usertoken.UserID, ubID)
	val, err := middleware.GetKey(followKey)
	if (err == nil && action_type == "1" && val == "1") || (err == nil && action_type == "2" && val == "0") {
		c.JSON(http.StatusOK, util.FollowActionErr)
		return
	}
	userA, _ := dao.GetUserByID(usertoken.UserID)                 // 查询数据库
	userB, _ := dao.GetUserByID(ubID)                             // 查询数据库
	err = middleware.UpdateRedisFollow(action_type, userA, userB) // 更新redis
	if err != nil {                                               // 关注操作失败
		fmt.Println(err)
		c.JSON(http.StatusOK, util.FollowActionErr)
		return
	}
	// 返回
	c.JSON(http.StatusOK, util.Success)
}

//关注列表函数，路由
func FollowList(c *gin.Context) {
	// 经过jwt中间件
	usertoken := c.MustGet("usertoken").(middleware.UserToken)
	// 查询之前，将Redis数据手动更新到Mysql
	SaveRedisToMySQL()
	//调用Dao层 查询关注者
	users, err := dao.GetFollow(usertoken.UserID)
	if err != nil {
		//返回关注列表和状态码
		c.JSON(http.StatusOK, FollowListResponse{
			Response: util.FollowListErr,
			UserList: []model.User{},
		})
	}
	// 查询redis，更新 用户的关注状态
	for i := range users {
		UpdateUserFollowStatus(usertoken.UserID, &users[i])
	}
	//返回关注列表和状态码
	c.JSON(http.StatusOK, FollowListResponse{
		Response: util.Success,
		UserList: users,
	})
}

//粉丝列表函数，路由
func FollowerList(c *gin.Context) {
	// 经过jwt中间件
	usertoken := c.MustGet("usertoken").(middleware.UserToken)
	// 查询之前，将Redis数据手动更新到Mysql
	SaveRedisToMySQL()
	//调用Dao层 查询粉丝
	users, err := dao.GetFollower(usertoken.UserID)
	if err != nil {
		//返回粉丝列表和状态码
		c.JSON(http.StatusOK, FollowListResponse{
			Response: util.FollowerListErr,
			UserList: []model.User{},
		})
	}
	// 查询redis，更新 用户的关注状态
	for i := range users {
		UpdateUserFollowStatus(usertoken.UserID, &users[i])
	}
	//返回粉丝列表和状态码
	c.JSON(http.StatusOK, FollowListResponse{
		Response: util.Success,
		UserList: users,
	})
}

// 更新 关注列表或者粉丝列表的用户的关注状态
func UpdateUserFollowStatus(userID int64, to_user *model.User) {
	redisKey := fmt.Sprintf("Fol:%d:%d", userID, to_user.Id)
	// 查询redis里面是否存在关注数
	if val, err := middleware.GetKey(redisKey); err != nil { //  若redis里面有数据，则以之为准
		if val == "1" {
			to_user.IsFollow = true
		} else {
			to_user.IsFollow = false
		}
	}
}
