package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yhmain/5th-simple-tiktok/model"
)

var ctx = context.Background()
var redisDB *redis.Client // 将redis的连接写成全局变量

// redis存放的数据种类  Fav: 与 Com: 都是前缀，以作区分
// 点赞状态->  	Key:  Fav:uid:vid			value: 1 或 0		会直接 Upset 到数据库，点赞表本身主键是uid:vid
// 点赞数量->  	Key:  FavCnt:vid			value: 数量			会更新tk_video表的favorite_count字段
// 评论信息->	Key:  ComAdd:cid  ComDel:cid value: 结构体序列化值	 会直接 Insert 到数据库，每生成一条评论就会有新comment_id
// 评论数量->	Key:  ComCnt:vid			value: 数量			会更新tk_video表的comment_count字段

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

// 查询key是否存在
func ExistKey(redisKey string) bool {
	res, _ := redisDB.Exists(ctx, redisKey).Result() // res代表存在的key数量
	return res > 0
}

// 获取Key对应的内容
func GetKey(redisKey string) (string, error) {
	val, err := redisDB.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		fmt.Println("key不存在")
		return "", err
	} else if err != nil {
		fmt.Println("获取key出现异常")
		return "", err
	} else {
		return val, nil
	}
}

// 设置 Key与值
func SetKey(redisKey, value string) (string, error) {
	return redisDB.Set(ctx, redisKey, value, 2*time.Minute).Result()
}

// 点赞 && 取消点赞 功能		key的过期时间设置为2min，0代表不会过期
func UpdateRedisLike(redisKey, vID, acType string) error {
	// 什么情况不用更改状态？ Key命中且当前操作类型与value代表含义一样
	val, err := GetKey(redisKey)
	if err == nil && acType == "1" && val == "1" { // 先获取当前 点赞的状态，防止用户多次点击造成的数据误加
		return errors.New("无效赞操作！")
	}
	if err == nil && acType == "2" && val == "0" {
		return errors.New("无效赞操作！")
	}
	if acType == "1" {
		redisDB.Set(ctx, redisKey, "1", 2*time.Minute) // 状态变为1
		redisDB.Incr(ctx, vID)                         // 数量加1
		fmt.Println("点赞成功!", "Redis key: ", redisKey)
	} else {
		redisDB.Set(ctx, redisKey, "0", 2*time.Minute) // 状态变为1
		redisDB.Decr(ctx, vID)                         // 数量减1
		fmt.Println("取消点赞成功!", "Redis key: ", redisKey)
	}
	return nil
}

// 新增评论与删除评论  需要保证结构体必须 有评论ID和视频ID
func UpdateRedisComment(acType string, comment model.Comment, comCount int64) {
	// 需要两个key，分别代表评论信息(ComAdd:cid  与  ComDel:cid)和评论数量(ComCnt:vid)
	commentStr, _ := json.Marshal(comment) // 结构体序列化
	redisAddCom := fmt.Sprintf("ComAdd:%d", comment.Id)
	redisComCnt := fmt.Sprintf("ComCnt:%d", comment.VideoID)
	if !ExistKey(redisComCnt) { // 首先将评论数插入redis
		redisDB.Set(ctx, redisComCnt, comCount, 2*time.Minute) // 加入到redsi中
	}
	if acType == "1" {
		// 添加评论
		redisDB.Set(ctx, redisAddCom, commentStr, 2*time.Minute) // 加入到redsi中
		redisDB.Incr(ctx, redisComCnt)                           // 评论数量加1
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
			redisDB.Set(ctx, redisDelCom, commentStr, 2*time.Minute) // 加入到redsi中
			redisDB.Decr(ctx, redisComCnt)                           // 评论数量减1
			fmt.Println("删除评论成功!", "Redis key: ", redisDelCom)
		}
	}
}

// 提取redis中的赞数据，包括：赞信息，赞数量
func GetRedisLike() ([]model.Like, map[int64]int64) {
	resLike, _ := getLikeByPattern("Fav:")  // 赞数据
	resCount, _ := getRedisCount("FavCnt:") //视频的赞数量
	return resLike, resCount
}

// 返回 新增评论和删除评论的ID
func GetRedisComment() ([]model.Comment, []model.Comment, map[int64]int64) {
	var patAdd = "ComAdd:"
	var patDel = "ComDel:"
	resAdd, _ := getCommentByPattern(patAdd) // 新增评论
	resDel, _ := getCommentByPattern(patDel) // 删除评论
	resCnt, _ := getRedisCount("ComCnt:")    // 视频的评论数量
	return resAdd, resDel, resCnt
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
		res := parseLikeByPattern(ks, pattern) // 解析key为  uid:vid
		likes = append(likes, res...)          // 合并到结果集中
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

// 解析 点赞类型的Key   key形式为： Fav:uid:vid  value为状态1或0
func parseLikeByPattern(ks []string, pat string) []model.Like {
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
