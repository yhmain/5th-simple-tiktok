package dao

import (
	"errors"
	"fmt"

	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 插入新的关注数据
func InsertNewFollow(follow *model.Follow) error {
	result := middleware.GetMySQLClient().Create(follow)
	return result.Error
}

// 按照 用户ID查询他关注的人
// select * form tk_user where uid in (select ubid from tk_follow where uaid=? and is_follow=1)
func GetFollow(userID int64) ([]model.User, error) {
	var users []model.User
	db := middleware.GetMySQLClient()
	err := db.Where("uid in (?)", db.Table("tk_follow").Where("uaid= ? and is_follow=1", userID).Select("ubid")).Find(&users).Error
	return users, err
}

// 按照用户ID查询他的粉丝
// select * form tk_user where uid in (select uaid from tk_follow where ubid=? and is_follow=1)
func GetFollower(userID int64) ([]model.User, error) {
	var users []model.User
	db := middleware.GetMySQLClient()
	err := db.Where("uid in (?)", db.Table("tk_follow").Where("ubid= ? and is_follow=1", userID).Select("uaid")).Find(&users).Error
	return users, err
}

// 更新关注表tk_follow，更新用户的 关注数和粉丝数
func SaveFollows(follows []model.Follow, followCount map[string]map[string]interface{}) error {
	db := middleware.GetMySQLClient()
	// 事务
	err := db.Transaction(func(tx *gorm.DB) error {
		// 先插入数据到 tk_follow表
		if len(follows) > 0 {
			result := db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "follow_id"}},            // key colume
				DoUpdates: clause.AssignmentColumns([]string{"is_follow"}), // column needed to be updated
			}).Create(&follows)
			if result.Error != nil {
				return errors.New("回滚：更新关注表")
			}
		}
		// 再更新 关注和粉丝 数量
		for k, v := range followCount { // k是用户ID，v是字典{关注数：xx，粉丝数：xx}
			err := UpdateUserFollow(k, v) //	更新数据库
			if err != nil {
				fmt.Println(err)
				return errors.New("回滚：更新用户关注数和粉丝数")
			}
		}
		return nil
	})
	return err
}

// 判断用户A是否关注用户B
// select count(*) from tk_follow where uaid=? and ubid=? and is_follow=1;
func IsAFollowB(uaid, ubid int64) bool {
	var count int64
	middleware.GetMySQLClient().Model(&model.Follow{}).Where("uaid=? and ubid=? and is_follow=1", uaid, ubid).Count(&count)
	return count > 0
}
