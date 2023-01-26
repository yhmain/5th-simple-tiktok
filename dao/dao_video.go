package dao

import (
	"github.com/yhmain/5th-simple-tiktok/model"
)

const (
	FEED_VIDEOS_NUM = 30 // 单次最多查询的数量为30条
)

//查询所有 <=latestTime 的 Video
func GetVideosByTime(latestTime int64) []model.Video {
	// 坑：preload里不是对应的表的名字，而是结构体中的字段名字！！！
	var videos []model.Video
	MyDB.Where("created_time<=?", latestTime).Order("created_time desc").Limit(FEED_VIDEOS_NUM).Preload("User").Find(&videos)
	return videos
}

//插入新发布的视频
func InsertNewVideo(video *model.Video) error {
	result := MyDB.Create(video)
	return result.Error
}

//获取某用户发布的所有视频
func GetVideosByUserID(UserID int64) []model.Video {
	var videos []model.Video
	MyDB.Where("uid=?", UserID).Preload("User").Find(&videos)
	return videos
}
