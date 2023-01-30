package middleware

import (
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/yhmain/5th-simple-tiktok/model"
)

// 测试 获取Redis连接
func TestGetRedisClient(t *testing.T) { // 测试ping，获取redis连接
	_, err := GetRedisClient()
	var expectedResult error
	assert.Equal(t, expectedResult, err)
}

// 测试解析redis的赞数据
func TestGetRedisLike(t *testing.T) {
	UpdateRedisLike("Fav:2:1618519735576563712", "1618519735576563712", "1")
	likes, err := GetRedisLike()
	assert.Equal(t, nil, err)
	expectedResult := []model.Like{
		{
			Id:         "2:1618519735576563712",
			UserID:     int64(2),
			VideoID:    int64(1618519735576563712),
			IsFavorite: true,
		},
	}
	assert.Equal(t, expectedResult, likes)
}

// 测试 评论的lua脚本
func TestLuaComment(t *testing.T) {
	// likeCntKey := fmt.Sprintf("FavCnt:%s", vID)
	// luaCommentScript := redis.NewScript(`
	// 	local acType = KEYS[1]
	// 	local likeKey = KEYS[2]
	// 	local likeCntKey = KEYS[3]
	// 	if (acType == "1")
	// 	then
	// 		-- 点赞操作
	// 		redis.call("SET", likeKey, 1, "EX", 120)		-- 状态变为1, 120表示120秒
	// 		redis.call("INCR", likeCntKey)					--赞数量+1
	// 	else
	// 		-- 取消点赞的操作
	// 		redis.call("SET", likeKey, 0, "EX", 120)		-- 状态变为0, 120表示120秒
	// 		redis.call("DECR", likeCntKey)					-- 赞数量-1
	// 	end
	// 	return 0
	// `)
	// n, err := luaLikeScript.Run(ctx, redisDB, []string{acType, likeKey, likeCntKey}).Result() //执行lua脚本
	// if err != nil {
	// 	fmt.Println("评论的lua脚本执行出现异常：", n, err)
	// 	return
	// }
}

// 测试关注操作
func TestUpdateRedisFollow(t *testing.T) {
	usera := model.User{
		Id:            1001,
		FollowCount:   0,
		FollowerCount: 2,
	}
	userb := model.User{
		Id:            2,
		FollowCount:   2,
		FollowerCount: 0,
	}
	UpdateRedisFollow("1", usera, userb)
	followKey := fmt.Sprintf("Fol:%d:%d", usera.Id, userb.Id)
	useraKey := fmt.Sprintf("FolCnt:%d", usera.Id)
	userbKey := fmt.Sprintf("FolCnt:%d", userb.Id)
	val, _ := GetKey(followKey)
	redisDB, _ := GetRedisClient()
	v1, _ := redisDB.HGet(ctx, useraKey, "FollowCount").Result()
	v2, _ := redisDB.HGet(ctx, userbKey, "FollowerCount").Result()
	assert.Equal(t, "1", val)                                     // 状态为1
	assert.Equal(t, fmt.Sprintf("%d", usera.FollowCount+1), v1)   // 关注数+1
	assert.Equal(t, fmt.Sprintf("%d", userb.FollowerCount+1), v2) //粉丝数+1
}

// 测试取关操作
func TestUpdateRedisFollow002(t *testing.T) {
	usera := model.User{
		Id:            1001,
		FollowCount:   0,
		FollowerCount: 2,
	}
	userb := model.User{
		Id:            2,
		FollowCount:   2,
		FollowerCount: 0,
	}
	UpdateRedisFollow("2", usera, userb)
	followKey := fmt.Sprintf("Fol:%d:%d", usera.Id, userb.Id)
	useraKey := fmt.Sprintf("FolCnt:%d", usera.Id)
	userbKey := fmt.Sprintf("FolCnt:%d", userb.Id)
	val, _ := GetKey(followKey)
	redisDB, _ := GetRedisClient()
	v1, _ := redisDB.HGet(ctx, useraKey, "FollowCount").Result()
	v2, _ := redisDB.HGet(ctx, userbKey, "FollowerCount").Result()
	assert.Equal(t, "0", val)                                     // 状态为1
	assert.Equal(t, fmt.Sprintf("%d", usera.FollowCount-1), v1)   // 关注数-1
	assert.Equal(t, fmt.Sprintf("%d", userb.FollowerCount-1), v2) //粉丝数-1
}

// 测试lua脚本写法
func TestLuaSCript(t *testing.T) {
	// if acType == "1" {
	// 	// 关注操作，a关注b
	// 	redisDB.Set(ctx, followKey, "1", keyTTL)           // 状态变为1
	// 	redisDB.HIncrBy(ctx, useraKey, "FollowCount", 1)   // a用户的关注数量加1
	// 	redisDB.HIncrBy(ctx, userbKey, "FollowerCount", 1) // b用户的粉丝数量加1
	// } else {
	// 	// 取消操作，a取关b
	// 	redisDB.Set(ctx, followKey, "0", keyTTL)            // 状态变为0
	// 	redisDB.HIncrBy(ctx, useraKey, "FollowCount", -1)   // a用户的关注数量-1
	// 	redisDB.HIncrBy(ctx, userbKey, "FollowerCount", -1) // b用户的粉丝数量-1
	// }
	luaFollowScript := redis.NewScript(`
		local acType = KEYS[1]
		local followKey = KEYS[2]
		local useraKey = KEYS[3]
		local userbKey = KEYS[4]
		if (acType == "1")
		then 
			-- 关注操作, a关注b了
			redis.call("SET", followKey, 1, "EX", 120)		-- 状态变为1
			redis.call("HINCRBY", useraKey, "FollowCount", 1)--a用户的关注数量+1
			redis.call("HINCRBY", userbKey, "FollowerCount", 1)--b用户的粉丝数量+1
		else 
			-- 取关操作, a取关b
			redis.call("SET", followKey, 0, "EX", 120)		-- 状态变为1
			redis.call("HINCRBY", useraKey, "FollowCount", -1)--a用户的关注数量-1
			redis.call("HINCRBY", userbKey, "FollowerCount", -1)--b用户的粉丝数量-1
		end
		return 0
	`)
	redisDB, _ := GetRedisClient()                                                                // 获取连接
	n, err := luaFollowScript.Run(ctx, redisDB, []string{"stock"}, []interface{}{+1}...).Result() // 传入4个参数
	if err != nil {
		fmt.Println("有异常：", n, err)
		return
	}
	fmt.Println("无异常：", n, err)
}
