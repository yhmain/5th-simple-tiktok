package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yhmain/5th-simple-tiktok/model"
	"gorm.io/gorm"
)

//插入新用户
func TestInserUser(t *testing.T) {
	err := InsertUser(&model.User{Id: 1001, UserName: "user1001", Password: "7777909"})
	var expectedResult error
	assert.Equal(t, expectedResult, err)
}

//该测试用例期待结果为空
func TestGetUserByName100(t *testing.T) {
	user, err := GetUserByName("wrong_user_name")
	expectedResult := model.User{Id: 0, UserName: "", Password: "", FollowCount: 0, FollowerCount: 0}
	assert.Equal(t, expectedResult, user)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

//该测试用例期待结果为 id为2的用户
func TestGetUserByName200(t *testing.T) {
	user, err := GetUserByName("admin")
	expectedResult := model.User{Id: 2, UserName: "admin", Password: "123456", FollowCount: 20, FollowerCount: 20}
	assert.Equal(t, expectedResult, user)
	assert.Equal(t, nil, err)
}

//测试按照用户名和密码搜索用户，查找成功
func TestGetUserByNamePwd100(t *testing.T) {
	user, err := GetUserByNamePwd("admin", "123456")
	expectedResult := model.User{Id: 2, UserName: "admin", Password: "123456", FollowCount: 20, FollowerCount: 20}
	assert.Equal(t, expectedResult, user)
	assert.Equal(t, nil, err)
}

//测试按照用户名和密码搜索用户，查找失败！
func TestGetUserByNamePwd200(t *testing.T) {
	user, err := GetUserByNamePwd("zhanglei", "123456")
	expectedResult := model.User{Id: 0, UserName: "", Password: "", FollowCount: 0, FollowerCount: 0}
	assert.Equal(t, expectedResult, user)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

//该测试用例期待结果为 id为2的用户
func TestGetUserByID001(t *testing.T) {
	user, err := GetUserByID(2)
	expectedResult := model.User{Id: 2, UserName: "admin", Password: "123456", FollowCount: 20, FollowerCount: 20}
	assert.Equal(t, expectedResult, user)
	assert.Equal(t, nil, err)
}

//该测试用例  不存在的用户ID
func TestGetUserByID002(t *testing.T) {
	_, err := GetUserByID(int64(2001))
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
