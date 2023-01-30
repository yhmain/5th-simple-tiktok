package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yhmain/5th-simple-tiktok/dao"
	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
	"github.com/yhmain/5th-simple-tiktok/util"
)

// 好友列表的响应体
type FriendListResponse struct {
	util.Response
	FriendList []model.Friend `json:"user_list"`
}

//好友列表函数，路由
func FriendList(c *gin.Context) {
	// 经过jwt中间件
	usertoken := c.MustGet("usertoken").(middleware.UserToken)
	// 查询数据库
	friends, err := dao.GetFriendsByID(usertoken.UserID)
	if err != nil {
		c.JSON(http.StatusOK, FriendListResponse{
			Response:   util.FriendListErr,
			FriendList: []model.Friend{},
		})
	}
	// 需要更新 User 信息
	// TODO...
	// 返回好友列表和状态码
	c.JSON(http.StatusOK, FriendListResponse{
		Response:   util.Success,
		FriendList: friends,
	})
}
