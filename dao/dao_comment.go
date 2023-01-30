package dao

import (
	"errors"

	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
	"gorm.io/gorm"
)

// 插入评论, 会导致视频的评论数目加1
func InsertNewComment(comment *model.Comment) error {
	result := middleware.GetMySQLClient().Create(comment)
	return result.Error
}

// 删除评论，会导致视频的评论数目减1
func DeleteComment(cID int64) error {
	result := middleware.GetMySQLClient().Where("cid=?", cID).Delete(&model.Comment{})
	return result.Error
}

// 保存评论的变化
func SaveComments(comAdd []model.Comment, comDel []model.Comment, comCount map[int64]int64) error {
	db := middleware.GetMySQLClient()
	// 事务
	err := db.Transaction(func(tx *gorm.DB) error {
		// 添加新评论
		if len(comAdd) > 0 {
			db.Create(&comAdd)
		}
		// 删除本该被删除的评论
		if len(comDel) > 0 {
			db.Delete(&comDel)
		}
		// 更新视频的评论数量
		for k, v := range comCount {
			err := UpdateVideoComment(k, v)
			if err != nil {
				return errors.New("回滚：更新视频评论数")
			}
		}
		return nil
	})
	return err
}

// 根据视频ID查询 里面的评论
func GetCommentsByVid(vID int64) ([]model.Comment, error) {
	var comments []model.Comment
	result := middleware.GetMySQLClient().Where("vid=?", vID).Preload("User").Preload("Video").Find(&comments)
	return comments, result.Error
}
