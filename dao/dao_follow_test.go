package dao

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yhmain/5th-simple-tiktok/model"
)

//插入新关注关系
func TestInsertNewFollow001(t *testing.T) {
	err := InsertNewFollow(&model.Follow{Id: "2:1618202276705341440", UserAID: 2, UserBID: 1618202276705341440, IsFollow: true})
	var expectedResult error
	assert.Equal(t, expectedResult, err)
}

func TestInsertNewFollow002(t *testing.T) {
	err := InsertNewFollow(&model.Follow{Id: "1000:1001", UserAID: 1000, UserBID: 1001, IsFollow: true})
	var expectedResult error
	assert.Equal(t, expectedResult, err)
}

// 测试 查找关注的人
func TestGetFollow(t *testing.T) {
	res, err := GetFollow(int64(2))
	fmt.Println(res)
	assert.Equal(t, nil, err)
}

// 测试 查找粉丝
func TestGetFollower(t *testing.T) {
	res, err := GetFollower(int64(1001))
	fmt.Println(res)
	assert.Equal(t, nil, err)
}
