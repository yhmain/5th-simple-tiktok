package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yhmain/5th-simple-tiktok/dao"
	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
	"github.com/yhmain/5th-simple-tiktok/util"
)

// 消息列表响应体
type MessageListResponse struct {
	util.Response
	MessageList []model.Message `json:"message_list"`
}

//发送消息函数，路由
func MessageChat(c *gin.Context) {
	// 经过jwt中间件
	usertoken := c.MustGet("usertoken").(middleware.UserToken)
	// 获取前端参数
	to_user_id := c.Query("to_user_id")             //接收方用户ID
	action_type := c.Query("action_type")           //操作类型
	content := c.Query("content")                   // 消息内容
	ubid, _ := strconv.ParseInt(to_user_id, 10, 64) //string转化成int64
	fmt.Println(to_user_id, action_type, content)
	if action_type == "1" {
		// 1. 发送消息
		mID := util.GenID() // 生成消息ID
		message := model.Message{
			Id:          mID,
			UserAID:     usertoken.UserID,
			UserBID:     ubid,
			Content:     content,
			CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
		}
		// 插入新消息
		if err := dao.InsertNewMessage(&message); err != nil {
			c.JSON(http.StatusOK, util.MessageActionErr)
			return
		}
	}
	// 返回
	c.JSON(http.StatusOK, util.Success)
}

//聊天记录函数，路由
func MessageAction(c *gin.Context) {
	// 经过jwt中间件
	usertoken := c.MustGet("usertoken").(middleware.UserToken)
	// 获取前端参数
	to_user_id := c.Query("to_user_id")             //接收方用户ID
	ubid, _ := strconv.ParseInt(to_user_id, 10, 64) //string转化成int64
	messages, err := dao.GetMessages(usertoken.UserID, ubid)
	if err != nil {
		c.JSON(http.StatusOK, MessageListResponse{
			Response:    util.MessageListErr,
			MessageList: []model.Message{},
		})
		return
	}
	// 返回成功
	c.JSON(http.StatusOK, MessageListResponse{
		Response:    util.Success,
		MessageList: messages,
	})
}
