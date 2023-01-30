package dao

import (
	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
)

// 新增好友 支持批量新增
func InsertNewFriend(friends interface{}) error {
	result := middleware.GetMySQLClient().Create(friends)
	return result.Error
}

// 查询 用户的好友列表
func GetFriendsByID(userID int64) ([]model.Friend, error) {
	var friends []model.Friend
	result := middleware.GetMySQLClient().Where("(uaid=? or ubid=?) and is_friend=1", userID, userID).Find(&friends)
	return friends, result.Error
}
