package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yhmain/5th-simple-tiktok/model"
)

//插入新用户
func TestInserUser(t *testing.T) {
	err := InsertUser(&model.User{Id: 1001, UserName: "user1001", Password: "7777909"})
	var expectedResult error
	assert.Equal(t, expectedResult, err)
}

//该测试用例期待结果为空
func TestGetUserByName100(t *testing.T) {
	user := GetUserByName("wrong_user_name")
	expectedResult := model.User{Id: 0, UserName: "", Password: "", FollowCount: 0, FollowerCount: 0}
	assert.Equal(t, expectedResult, user)
}

//该测试用例期待结果为 id为2的用户
func TestGetUserByName200(t *testing.T) {
	user := GetUserByName("admin")
	expectedResult := model.User{Id: 2, UserName: "admin", Password: "123456", FollowCount: 20, FollowerCount: 20}
	assert.Equal(t, expectedResult, user)
}

//测试按照用户名和密码搜索用户，查找成功
func TestGetUserByNamePwd100(t *testing.T) {
	user := GetUserByNamePwd("admin", "123456")
	expectedResult := model.User{Id: 2, UserName: "admin", Password: "123456", FollowCount: 20, FollowerCount: 20}
	assert.Equal(t, expectedResult, user)
}

//测试按照用户名和密码搜索用户，查找失败！
func TestGetUserByNamePwd200(t *testing.T) {
	user := GetUserByNamePwd("zhanglei", "123456")
	expectedResult := model.User{Id: 0, UserName: "", Password: "", FollowCount: 0, FollowerCount: 0}
	assert.Equal(t, expectedResult, user)
}

//该测试用例期待结果为 id为2的用户
func TestGetUserByID(t *testing.T) {
	user := GetUserByID(2)
	expectedResult := model.User{Id: 2, UserName: "admin", Password: "123456", FollowCount: 20, FollowerCount: 20}
	assert.Equal(t, expectedResult, user)
}
