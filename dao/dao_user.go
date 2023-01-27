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
func GetUserByName(name string) model.User {
	var user model.User
	middleware.GetMySQLClient().Where("user_name=?", name).Find(&user)
	return user
}

//查找用户名、密码是否正确, nil空则表示不存在
func GetUserByNamePwd(name, pwd string) model.User {
	var user model.User
	middleware.GetMySQLClient().Where("user_name=? AND password=?", name, pwd).Find(&user)
	return user
}

//根据ID查找用户信息
func GetUserByID(uid int64) model.User {
	var user model.User
	middleware.GetMySQLClient().Where("uid=?", uid).Find(&user)
	return user
}
