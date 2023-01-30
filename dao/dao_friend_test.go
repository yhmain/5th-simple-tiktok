package dao

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
)

// 测试插入 好友数据
func TestInsertNewFriend(t *testing.T) {
	db := middleware.GetMySQLClient()
	selectUserID := int64(2) // 为这个人添加所有的用户作为好友
	var users []model.User
	db.Where("uid <> ?", selectUserID).Find(&users)
	fmt.Println("查询出来的全部用户：", users)
	var friends []model.Friend
	for i := range users { // 构造朋友结构体
		ua := selectUserID
		ub := users[i].Id
		if ua > ub {
			ua, ub = ub, ua
		}
		friend := model.Friend{
			Id:       fmt.Sprintf("%d:%d", ua, ub),
			UserAID:  ua,
			UserBID:  ub,
			IsFriend: true,
		}
		friends = append(friends, friend)
	}
	err := InsertNewFriend(&friends)
	assert.Equal(t, nil, err)
}

// 测试 好友列表查询功能
func TestGetFriendsByID(t *testing.T) {
	friends, err := GetFriendsByID(int64(2))
	fmt.Println("好友数量：", len(friends))
	fmt.Println(friends)
	assert.Equal(t, nil, err)
}
