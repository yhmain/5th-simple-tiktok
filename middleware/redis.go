package middleware

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yhmain/5th-simple-tiktok/model"
)

var ctx = context.Background()
var redisDB *redis.Client // 将redis的连接写成全局变量

func init() {
	// init函数内进行初始化
	redisDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// 测试redis连接
	_, err := redisDB.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Redis连接出错: ", err)
	}
	redisDB.FlushAll(ctx) // 此处清空缓存
}

// 获取客户端连接
func GetRedisClient() (*redis.Client, error) {
	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := redisDB.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return redisDB, nil
}

// 供外层调用：清空redis缓存
func ClearRedis() {
	redisDB.FlushAll(ctx) // 此处清空缓存
}

// 解析Redis的所有key，保存到Like结构体
func ParseRedisKeys() ([]model.Like, error) {
	// 保存到数据库为批量操作
	// 假设该台redis服务器所有Key都是uid:vid形式
	ks, err := redisDB.Keys(ctx, "Fav:*").Result() // 获取Redis的所有key
	if err != nil {
		fmt.Println("获取Redis的Keys出错: ", err)
		return nil, err
	}
	likes := make([]model.Like, 0)
	for _, k := range ks {
		realID := k[4:]
		ss := strings.Split(realID, ":") // uid:vid
		uid, _ := strconv.ParseInt(ss[0], 10, 64)
		vid, _ := strconv.ParseInt(ss[1], 10, 64)
		value, _ := redisDB.Get(ctx, k).Result() // 获取该Key对应的值
		var val = false
		if value == "1" { // 1表示点赞
			val = true
		}
		lk := model.Like{ // 构造 点赞结构体
			Id:         realID,
			UserID:     uid,
			VideoID:    vid,
			IsFavorite: val,
		}
		likes = append(likes, lk)
	}
	return likes, nil
}

// 点赞 && 取消点赞 功能
// Redis存储的数据：
// 1. 用户点赞，key为uid:vid的组合形式,value 为 0或1
// key的过期时间设置为2min，0代表不会过期
func UpdateRedisLike(redisKey, vID, acType string) (string, error) {
	val, err := GetKey(redisKey) // 先获取当前 点赞的状态，防止用户多次点击造成的数据误加
	if err == nil {              // 若redis中有该key
		if acType == "1" && val == "0" { // 1 为点赞操作
			fmt.Println("点赞成功!", "Redis key: ", redisKey)
			return redisDB.Set(ctx, redisKey, "1", 2*time.Minute).Result()
		} else if acType == "2" && val == "1" { // 2为取消点赞的操作
			fmt.Println("取消点赞成功!", "Redis key: ", redisKey)
			return redisDB.Set(ctx, redisKey, "0", 2*time.Minute).Result()
		}
		fmt.Println("无效点赞操作!", "Redis key: ", redisKey, "赞状态：", val)
		return "", nil
	} else {
		if acType == "1" { // 1 为点赞操作
			fmt.Println("点赞成功!", "Redis key: ", redisKey)
			return redisDB.Set(ctx, redisKey, "1", 2*time.Minute).Result()
		} else { // 2为取消点赞的操作
			fmt.Println("取消点赞成功!", "Redis key: ", redisKey)
			return redisDB.Set(ctx, redisKey, "0", 2*time.Minute).Result()
		}
	}
}

// 查询key是否存在于redis
func GetKey(redisKey string) (string, error) {
	val, err := redisDB.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		fmt.Println("Redis key does not exist.", redisKey)
		return "Error", err
	} else if err != nil {
		fmt.Println("Redis error")
		return "Error", err
	} else {
		return val, nil
	}
}
