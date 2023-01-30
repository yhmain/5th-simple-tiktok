package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yhmain/5th-simple-tiktok/model"
)

// 测试插入新消息
func TestInsertNewMessage(t *testing.T) {
	message := model.Message{
		Id:          int64(1000),
		UserAID:     int64(2),
		UserBID:     int64(1002),
		Content:     "这是admin给user发送的消息",
		CreatedTime: "2023-01-29 20:57:46",
	}
	err := InsertNewMessage(&message)
	assert.Equal(t, nil, err)
}
