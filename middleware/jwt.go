package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/yhmain/5th-simple-tiktok/util"
)

//NOTE: 以下为基于jwt-go实现的token权限认证

// 用户Token结构体  JWT:Json Web Token
type UserToken struct {
	UserID   int64
	Name     string
	Password string
}

//用于生成Token的结构体
type UserClaims struct {
	UserToken
	jwt.RegisteredClaims //v4版本新增
}

var jwtKey = []byte("tiktok")             //定义Secret
const TokenExpireDuration = time.Hour * 2 //定义JWT的过期时间：2小时

//发放Token
func GenToken(userToken *UserToken) (string, error) {
	// 创建一个用户的声明，即初始化结构体 UserClaims
	c := UserClaims{
		UserToken: *userToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                          // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                          // 生效时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c) // 使用指定的签名方法（如HS256）创建签名对象
	return token.SignedString(jwtKey)                     // 使用指定的Secret签名并获得完整的编码后的字符串token
}

//解析Token
func ParseToken(tokenString string) (*jwt.Token, *UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}

// 自定义函数：JWTAuthUser 基于JWT认证的中间件
func JWTAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		//获取token，并解析
		token := GetParamPostOrGet(c, "token")
		// fmt.Println("获取到的Token: ", token)
		_, claims, err := ParseToken(token)
		if err != nil {
			fmt.Println("Token解析出错: ", err)
			c.JSON(http.StatusOK, util.ParseTokenErr) //token解析失败
			c.Abort()
			return
		}
		// 获取user_id，并与token解析出来的进行对比
		// 若存在user_id，则需要进行鉴权
		paramID := GetParamPostOrGet(c, "user_id")
		if paramID != "" && strconv.FormatInt(claims.UserID, 10) != paramID {
			c.JSON(http.StatusOK, util.WrongTokenErr) //token校验失败
			c.Abort()
			return
		}
		c.Set("usertoken", claims.UserToken)
		c.Next() // 执行后续的处理函数
	}
}

// 发现post和get有时候与给的不一致
func GetParamPostOrGet(c *gin.Context, param string) string {
	token := ""
	t1 := c.Query(param)
	t2 := c.PostForm(param)
	if t1 != "" {
		token = t1
	}
	if t2 != "" {
		token = t2
	}
	return token
}
