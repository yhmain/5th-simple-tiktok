package dao

import (
	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
)

// 插入新消息
func InsertNewMessage(message *model.Message) error {
	result := middleware.GetMySQLClient().Create(message)
	return result.Error
}

// 查询消息记录
func GetMessages(uaid int64, ubid int64) ([]model.Message, error) {
	var mes []model.Message
	result := middleware.GetMySQLClient().Where("uaid=? and ubid=?", uaid, ubid).Find(&mes)
	return mes, result.Error
}
