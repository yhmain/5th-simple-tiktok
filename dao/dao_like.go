package dao

import (
	"errors"

	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 将 redis中的赞数据持久化到Mysql中，有如下事务
// 1. 赞数据更新到 tk_like表
// 2. 更新对应视频的赞数量 tk_video表
func SaveLikes(likes []model.Like, likeCount map[int64]int64) error {
	db := middleware.GetMySQLClient()
	// 事务
	err := db.Transaction(func(tx *gorm.DB) error {
		if len(likes) > 0 { // 空切片不做处理
			// 先插入数据到 tk_like表
			result := db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "like_id"}},                // key colume
				DoUpdates: clause.AssignmentColumns([]string{"is_favorite"}), // column needed to be updated
			}).Create(&likes)
			if result.Error != nil {
				return errors.New("回滚：更新点赞表")
			}
		}
		// 再更新赞数量
		for k, v := range likeCount {
			err := UpdateVideoLike(k, v)
			if err != nil {
				return errors.New("回滚：更新视频点赞数")
			}
		}
		return nil
	})
	return err
}

// 查询用户是否喜欢该视频
func GetLikeByID(likeID string) model.Like {
	var like model.Like
	middleware.GetMySQLClient().Where("like_id=?", likeID).Preload("User").Preload("Video").Find(&like) // 只是想得到是否喜欢这个数据，故不考虑预加载
	return like
}

// 根据用户ID查询，他点赞的视频		为了，优化，可以建立索引(uid)
func GetUserLike(userID int64) []model.Like {
	var likes []model.Like
	middleware.GetMySQLClient().Where("uid=? and is_favorite=?", userID, 1).Preload("User").Preload("Video").Find(&likes)
	return likes
}
