package dao

import (
	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
)

//插入新用户，返回是否出错，为nil则表示插入成功
func InsertUser(user *model.User) error {
	result := middleware.GetMySQLClient().Create(user)
	return result.Error
}

//查找用户名是否已存在（按照规定，用户名是唯一的）, nil空则表示不存在
func GetUserByName(name string) (model.User, error) {
	var user model.User
	result := middleware.GetMySQLClient().Where("user_name=?", name).First(&user)
	return user, result.Error
}

//查找用户名、密码是否正确, nil空则表示不存在
func GetUserByNamePwd(name, pwd string) (model.User, error) {
	var user model.User
	result := middleware.GetMySQLClient().Where("user_name=? AND password=?", name, pwd).First(&user)
	return user, result.Error
}

//根据ID查找用户信息
func GetUserByID(uid int64) (model.User, error) {
	var user model.User
	result := middleware.GetMySQLClient().Where("uid=?", uid).First(&user)
	return user, result.Error
}

// 更新用户的关注数和粉丝数
func UpdateUserFollow(uid string, data map[string]interface{}) error {
	result := middleware.GetMySQLClient().Model(&model.User{}).Where("uid=?", uid).Updates(data)
	return result.Error
}
