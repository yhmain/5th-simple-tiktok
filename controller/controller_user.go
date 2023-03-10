package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yhmain/5th-simple-tiktok/dao"
	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
	"github.com/yhmain/5th-simple-tiktok/util"
)

//用户登录的响应体
type UserLoginResponse struct {
	util.Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

//用户信息响应体
type UserResponse struct {
	util.Response
	User model.User `json:"user"`
}

//用户注册函数，路由
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	//用户名和密码的长度不应该超过32
	if len(username) > 32 || len(password) > 32 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: util.TooLongInputErr,
			UserId:   -1,
			Token:    "",
		})
		return
	}
	//判断该用户名是否已存在，为了确保用户名是唯一的
	if _, err := dao.GetUserByName(username); err == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: util.NameAlreadyExistErr,
			UserId:   -1,
			Token:    "",
		})
	} else {
		//生成新用户的ID
		newUserID := util.GenID()
		//构造新的用户结构体
		newUser := model.User{
			Id:       newUserID,
			UserName: username,
			Password: password,
		}
		// 插入数据库
		if err := dao.InsertUser(&newUser); err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: util.InsertNewUserErr,
				UserId:   -1,
				Token:    "",
			})
			return
		}
		//生成新用户的鉴权token
		token, err := middleware.GenToken(&middleware.UserToken{
			UserID:   newUserID,
			Name:     username,
			Password: password,
		})
		if err != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: util.GenTokenFailedErr,
				UserId:   -1,
				Token:    "",
			})
			return
		}
		//最后，返回成功！！！
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: util.Success,
			UserId:   newUserID,
			Token:    token,
		})
	}
}

//用户登录函数，路由
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	//检测 用户名和密码是否正确
	if user, err := dao.GetUserByNamePwd(username, password); err != nil {
		//校验成功后，生成用户鉴权token
		token, err := middleware.GenToken(&middleware.UserToken{UserID: user.Id, Name: username, Password: password})
		if err != nil {
			//若生成token出错，则返回错误代码
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: util.GenTokenFailedErr,
				UserId:   -1,
				Token:    "",
			})
			return
		}

		//成功获取token
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: util.Success,
			UserId:   user.Id,
			Token:    token,
		})

	} else { //用户名或者密码错误
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: util.WrongPassword,
			UserId:   -1,
			Token:    "",
		})
	}
}

//输入为用户id和鉴权token，获取该用户信息
//注意：有中间件已处理
func UserInfo(c *gin.Context) {
	usertoken := c.MustGet("usertoken").(middleware.UserToken)

	//按照id查找用户信息
	user, err := dao.GetUserByID(usertoken.UserID)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: util.WrongUserID,
			User:     model.User{}, // 返回空
		})
		return
	}
	// 查询redis，更新登录用户的关注/粉丝数量
	UpdateUserFollow(&user)
	//此时代表成功！
	c.JSON(http.StatusOK, UserResponse{
		Response: util.Success,
		User:     user,
	})
}

// 更新用户的关注数和粉丝数
func UpdateUserFollow(user *model.User) {
	redisKey := fmt.Sprintf("FolCnt:%d", user.Id)
	// 查询redis里面是否存在关注数
	if val, err := middleware.HGetKey(redisKey, "FollowCount"); err != nil { //  若redis里面有数据，则以之为准
		cnt, _ := strconv.ParseInt(val, 10, 64)
		user.FollowCount = cnt // 更新
	}
	// 查询redis里面是否存在粉丝数
	if val, err := middleware.HGetKey(redisKey, "FollowCount"); err != nil { //  若redis里面有数据，则以之为准
		cnt, _ := strconv.ParseInt(val, 10, 64)
		user.FollowerCount = cnt // 更新
	}
}
