package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yhmain/5th-simple-tiktok/model"
)

var ctx = context.Background()
var redisDB *redis.Client    // 将redis的连接写成全局变量
var keyTTL = time.Minute * 2 // 假设redis的key过期时间是2分钟

// redis存放的数据类型  Fav: 与 Com: 都是前缀，以作区分
// 点赞状态->  	Key:  Fav:uid:vid			value: 1 或 0		会直接 Upset 到数据库，点赞表本身主键是uid:vid
// 点赞数量->  	Key:  FavCnt:vid			value: 数量			会更新tk_video表的favorite_count字段
// 评论信息->	Key:  ComAdd:cid  ComDel:cid value: 结构体序列化值	 会直接 Insert 到数据库，每生成一条评论就会有新comment_id
// 评论数量->	Key:  ComCnt:vid			value: 数量			会更新tk_video表的comment_count字段
// 关注信息->	Key:  Fol:uaid:ubid			value :1 或 0		会直接 Upset 到数据库，关注表本身主键是uaid:ubid
// 对于每个用户的关注/粉丝数量，有 Key: FolCnt:uid  field:"FollowCount"  field:"FollowerCount"，采取哈希结构

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
	redisDB.FlushAll(ctx) // 假设重新连接redis，清空缓存
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

// 查询key是否存在
func ExistKey(redisKey string) bool {
	res, _ := redisDB.Exists(ctx, redisKey).Result() // res代表存在的key数量
	return res > 0
}

// 获取Key对应的内容
func GetKey(redisKey string) (string, error) {
	val, err := redisDB.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		// fmt.Println("key不存在", redisKey)
		return "", err
	} else if err != nil {
		// fmt.Println("获取key出现异常", redisKey)
		return "", err
	} else {
		return val, nil
	}
}

// 设置 Key与值
func SetKey(redisKey, value string) (string, error) {
	return redisDB.Set(ctx, redisKey, value, keyTTL).Result()
}

// 点赞 && 取消点赞 功能		key的过期时间设置为2min，0代表不会过期
func UpdateRedisLike(likeKey, vID, acType string) error {
	likeCntKey := fmt.Sprintf("FavCnt:%s", vID)
	luaLikeScript := redis.NewScript(`
		local acType = KEYS[1]
		local likeKey = KEYS[2]
		local likeCntKey = KEYS[3]
		if (acType == "1")
		then
			-- 点赞操作
			redis.call("SET", likeKey, 1, "EX", 120)		-- 状态变为1, 120表示120秒
			redis.call("INCR", likeCntKey)					--赞数量+1
		else 
			-- 取消点赞的操作
			redis.call("SET", likeKey, 0, "EX", 120)		-- 状态变为0, 120表示120秒
			redis.call("DECR", likeCntKey)					-- 赞数量-1
		end
		return 0
	`)
	n, err := luaLikeScript.Run(ctx, redisDB, []string{acType, likeKey, likeCntKey}).Result() //执行lua脚本
	if err != nil {
		fmt.Println("点赞的lua脚本执行出现异常：", n, err)
		return err
	}
	return nil
}

// 新增评论与删除评论  需要保证结构体必须 有评论ID和视频ID
func UpdateRedisComment(acType string, comment model.Comment) {
	// 需要两个key，分别代表评论信息(ComAdd:cid  与  ComDel:cid)和评论数量(ComCnt:vid)
	commentStr, _ := json.Marshal(comment) // 结构体序列化
	redisAddCom := fmt.Sprintf("ComAdd:%d", comment.Id)
	redisComCnt := fmt.Sprintf("ComCnt:%d", comment.VideoID)
	if acType == "1" {
		// 添加评论
		redisDB.Set(ctx, redisAddCom, commentStr, keyTTL) // 加入到redsi中
		redisDB.Incr(ctx, redisComCnt)                    // 评论数量加1
		fmt.Println("发布评论成功!", "Redis key: ", redisAddCom)
	} else {
		// 删除评论
		if ExistKey(redisAddCom) {
			// 如果要删除的Key在  AddCom:cid 里面，则直接删除redis里面的Key
			redisDB.Del(ctx, redisAddCom)
			redisDB.Incr(ctx, redisComCnt) // 评论数量减1
			fmt.Println("删除评论成功!", "Redis key: ", redisAddCom)
		} else {
			// 否则添加到   DelCom:cid里面  后面会去删除数据库的内容
			redisDelCom := fmt.Sprintf("ComDel:%d", comment.Id)
			redisDB.Set(ctx, redisDelCom, commentStr, keyTTL) // 加入到redsi中
			redisDB.Decr(ctx, redisComCnt)                    // 评论数量减1
			fmt.Println("删除评论成功!", "Redis key: ", redisDelCom)
		}
	}
}

// 关注操作，关注与取关；传进来的两个用户模型必须携带用户ID、关注数、粉丝数
func UpdateRedisFollow(acType string, usera, userb model.User) error {
	followKey := fmt.Sprintf("Fol:%d:%d", usera.Id, userb.Id)
	useraKey := fmt.Sprintf("FolCnt:%d", usera.Id)
	userbKey := fmt.Sprintf("FolCnt:%d", userb.Id)
	if !ExistKey(useraKey) { // 将用户a的两个数据插入redis
		redisDB.HMSet(ctx, useraKey, map[string]interface{}{"FollowCount": usera.FollowCount, "FollowerCount": usera.FollowerCount})
	}
	if !ExistKey(userbKey) { // 将用户b的两个数据插入redis
		redisDB.HMSet(ctx, userbKey, map[string]interface{}{"FollowCount": userb.FollowCount, "FollowerCount": userb.FollowerCount})
	}
	luaFollowScript := redis.NewScript(`
		local acType = KEYS[1]
		local followKey = KEYS[2]
		local useraKey = KEYS[3]
		local userbKey = KEYS[4]
		if (acType == "1")
		then 
			-- 关注操作, a关注b了
			redis.call("SET", followKey, 1, "EX", 120)		-- 状态变为1, 120表示120秒
			redis.call("HINCRBY", useraKey, "FollowCount", 1)--a用户的关注数量+1
			redis.call("HINCRBY", userbKey, "FollowerCount", 1)--b用户的粉丝数量+1
		else 
			-- 取关操作, a取关b
			redis.call("SET", followKey, 0, "EX", 120)		-- 状态变为0, 120表示120秒
			redis.call("HINCRBY", useraKey, "FollowCount", -1)-- a用户的关注数量-1
			redis.call("HINCRBY", userbKey, "FollowerCount", -1)-- b用户的粉丝数量-1
		end
		return 0
	`)
	n, err := luaFollowScript.Run(ctx, redisDB, []string{acType, followKey, useraKey, userbKey}).Result() //执行lua脚本
	if err != nil {
		fmt.Println("关注的lua脚本执行出现异常：", n, err)
		return err
	}
	return nil
}

// 提取redis中的赞数据，包括：赞信息，赞数量
func GetRedisLike() ([]model.Like, map[int64]int64) {
	resLike, _ := getLikeByPattern("Fav:")  // 赞数据
	resCount, _ := getRedisCount("FavCnt:") //视频的赞数量
	return resLike, resCount
}

// 返回 新增评论和删除评论的ID
func GetRedisComment() ([]model.Comment, []model.Comment, map[int64]int64) {
	resAdd, _ := getCommentByPattern("ComAdd:") // 新增评论
	resDel, _ := getCommentByPattern("ComDel:") // 删除评论
	resCnt, _ := getRedisCount("ComCnt:")       // 视频的评论数量
	return resAdd, resDel, resCnt
}

// 返回 关注信息 和 用户ID:关注数和粉丝数
func GetRedisFollow() ([]model.Follow, map[string]map[string]interface{}) {
	resFollow, _ := getFollowByPattern("Fol:") // 关注信息
	resCount, _ := getFollowCount("FolCnt:")   // 关注数量与粉丝数量
	return resFollow, resCount
}

// 解析点赞状态，保存到Like结构体
func getLikeByPattern(pattern string) ([]model.Like, error) {
	// var pattern = "Fav:"       // 模糊查询的前缀
	var cursor uint64          // 定义游标
	var batchSize = int64(100) // 定义每次获取多少key
	var likes = []model.Like{} // 最终返回的结果
	for {
		var err error
		// 扫描所有key，每次100条，比Keys方法有优势，因为不会阻塞redis主线程
		ks, cursor, err := redisDB.Scan(ctx, cursor, pattern+"*", batchSize).Result()
		if err != nil {
			fmt.Println("获取Redis的Keys出错: ", err)
			return nil, err
		}
		res := parseLike(ks, pattern) // 解析key为  uid:vid
		likes = append(likes, res...) // 合并到结果集中
		// 考虑在此时, 删除这部分key
		for i := 0; i < len(ks); i++ {
			_, err = redisDB.Del(ctx, ks[i]).Result()
			if err != nil {
				fmt.Println("删除Redis的Keys出错: ", err)
				return nil, err
			}
		}
		if cursor == 0 {
			break
		}
	}
	return likes, nil
}

// 解析评论信息，保存到评论结构体
func getCommentByPattern(pattern string) ([]model.Comment, error) {
	var cursor uint64            // 定义游标
	var batchSize = int64(100)   // 定义每次获取多少key
	var coms = []model.Comment{} // 最终返回的结果
	for {
		var err error
		// 扫描所有key，每次100条，比Keys方法有优势，因为不会阻塞redis主线程
		ks, cursor, err := redisDB.Scan(ctx, cursor, pattern+"*", batchSize).Result()
		if err != nil {
			fmt.Println("获取Redis的Keys出错: ", err)
			return coms, err
		}
		// 解析 点赞数量与视频ID的映射关系   key形式为： ComAdd:cid  value为 序列化的字符串
		for _, k := range ks {
			value, _ := redisDB.Get(ctx, k).Result() // 获取该Key对应的值
			// 合并到结果集中
			var com model.Comment
			json.Unmarshal([]byte(value), &com)
			coms = append(coms, com)
		}
		// 考虑在此时, 删除这部分key
		for i := 0; i < len(ks); i++ {
			_, err = redisDB.Del(ctx, ks[i]).Result()
			if err != nil {
				fmt.Println("删除Redis的Keys出错: ", err)
				return coms, err
			}
		}
		if cursor == 0 {
			break
		}
	}
	return coms, nil
}

// 解析关注状态，保存到Follow结构体
func getFollowByPattern(pattern string) ([]model.Follow, error) {
	var cursor uint64              // 定义游标
	var batchSize = int64(100)     // 定义每次获取多少key
	var follows = []model.Follow{} // 最终返回的结果
	for {
		var err error
		// 扫描所有key，每次100条，比Keys方法有优势，因为不会阻塞redis主线程
		ks, cursor, err := redisDB.Scan(ctx, cursor, pattern+"*", batchSize).Result()
		if err != nil {
			fmt.Println("获取Redis的Keys出错: ", err)
			return nil, err
		}
		res := parseFollow(ks, pattern)   // 解析key为  uaid:ubid
		follows = append(follows, res...) // 合并到结果集中
		// 考虑在此时, 删除这部分key
		for i := 0; i < len(ks); i++ {
			_, err = redisDB.Del(ctx, ks[i]).Result()
			if err != nil {
				fmt.Println("删除Redis的Keys出错: ", err)
				return nil, err
			}
		}
		if cursor == 0 {
			break
		}
	}
	return follows, nil
}

// 解析，得到用户的关注数和粉丝数
func getFollowCount(pattern string) (map[string]map[string]interface{}, error) {
	var cursor uint64                                      // 定义游标
	var batchSize = int64(100)                             // 定义每次获取多少key
	var followCounts = map[string]map[string]interface{}{} // 最终返回的结果
	for {
		var err error
		// 扫描所有key，每次100条，比Keys方法有优势，因为不会阻塞redis主线程
		ks, cursor, err := redisDB.Scan(ctx, cursor, pattern+"*", batchSize).Result()
		if err != nil {
			fmt.Println("获取Redis的Keys出错: ", err)
			return followCounts, err
		}
		// 解析 点赞数量与视频ID的映射关系   key形式为： ComAdd:cid  value为 序列化的字符串
		for _, k := range ks {
			val1, _ := redisDB.HGet(ctx, k, "FollowCount").Result()   // 获取该Key对应的值
			val2, _ := redisDB.HGet(ctx, k, "FollowerCount").Result() // 获取该Key对应的值
			fCnt, _ := strconv.ParseInt(val1, 10, 64)
			ferCnt, _ := strconv.ParseInt(val2, 10, 64)
			// userID, _ := strconv.ParseInt(k[len(pattern):], 10, 64)
			fmt.Println("关注数：", k, fCnt, ferCnt)
			followCounts[k[len(pattern):]] = map[string]interface{}{"FollowCount": fCnt, "FollowerCount": ferCnt} // 构造成字典
		}
		// 考虑在此时, 删除这部分key
		for i := 0; i < len(ks); i++ {
			_, err = redisDB.Del(ctx, ks[i]).Result()
			if err != nil {
				fmt.Println("删除Redis的Keys出错: ", err)
				return followCounts, err
			}
		}
		if cursor == 0 {
			break
		}
	}
	return followCounts, nil
}

// 解析 点赞类型的Key   key形式为： Fav:uid:vid  value为状态1或0
func parseLike(ks []string, pat string) []model.Like {
	likes := make([]model.Like, 0)
	for _, k := range ks {
		realID := k[len(pat):]           // 因为 fav:  长度是4
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
	return likes
}

// 解析 关注类型的Key   key形式为： Fol:uaid:ubid  value为状态1或0
func parseFollow(ks []string, pat string) []model.Follow {
	follows := make([]model.Follow, 0)
	for _, k := range ks {
		realID := k[len(pat):]
		ss := strings.Split(realID, ":") // uaid:ubid
		uaid, _ := strconv.ParseInt(ss[0], 10, 64)
		ubid, _ := strconv.ParseInt(ss[1], 10, 64)
		value, _ := redisDB.Get(ctx, k).Result() // 获取该Key对应的值
		var val = false
		if value == "1" { // 1表示关注
			val = true
		}
		fol := model.Follow{
			Id:       realID,
			UserAID:  uaid,
			UserBID:  ubid,
			IsFollow: val,
		}
		follows = append(follows, fol)
	}
	return follows
}

// 解析点赞或评论的数量，保存到字典中
func getRedisCount(pattern string) (map[int64]int64, error) {
	// var pattern = "FavCnt:"           // 模糊查询的前缀
	var cursor uint64                 // 定义游标
	var batchSize = int64(100)        // 定义每次获取多少key
	var likeCount = map[int64]int64{} // 最终返回的结果
	for {
		var err error
		// 扫描所有key，每次100条，比Keys方法有优势，因为不会阻塞redis主线程
		ks, cursor, err := redisDB.Scan(ctx, cursor, pattern+"*", batchSize).Result()
		if err != nil {
			fmt.Println("获取Redis的Keys出错: ", err)
			return nil, err
		}
		// 解析 点赞数量与视频ID的映射关系   key形式为： FavCnt:vid  value为赞数量
		for _, k := range ks {
			realID := k[len(pattern):]
			value, _ := redisDB.Get(ctx, k).Result() // 获取该Key对应的值
			vid, _ := strconv.ParseInt(realID, 10, 64)
			cnt, _ := strconv.ParseInt(value, 10, 64)
			likeCount[vid] = cnt
		}
		// 考虑在此时, 删除这部分key
		for i := 0; i < len(ks); i++ {
			_, err = redisDB.Del(ctx, ks[i]).Result()
			if err != nil {
				fmt.Println("删除Redis的Keys出错: ", err)
				return nil, err
			}
		}
		if cursor == 0 {
			break
		}
	}
	return likeCount, nil
}
