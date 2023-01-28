package dao

import (
	"github.com/yhmain/5th-simple-tiktok/middleware"
	"github.com/yhmain/5th-simple-tiktok/model"
)

const (
	FEED_VIDEOS_NUM = 30 // 单次最多查询的数量为30条
)

//查询所有 <=latestTime 的 Video
func GetVideosByTime(latestTime int64) []model.Video {
	// 坑：preload里不是对应的表的名字，而是结构体中的字段名字！！！
	var videos []model.Video
	middleware.GetMySQLClient().Where("created_time<=?", latestTime).Order("created_time desc").Limit(FEED_VIDEOS_NUM).Preload("User").Find(&videos)
	return videos
}

// 根据一组视频ID 批量查询视频
func GetVideoByIDs(videoIDs []int64) []model.Video {
	var videos []model.Video
	middleware.GetMySQLClient().Where("vid in ?", videoIDs).Find(&videos)
	return videos
}

// 根据视频ID 更新视频的点赞数
func UpdateVideoLike(videoID, favCount int64) error {
	// update tk_video set favorite_count=favorite_count+{1} where vid={2}
	result := middleware.GetMySQLClient().Model(&model.Video{}).Where("vid=?", videoID).Update("favorite_count", favCount)
	return result.Error
}

// 根据视频ID 更新视频的评论数
func UpdateVideoComment(videoID, comCount int64) error {
	result := middleware.GetMySQLClient().Model(&model.Video{}).Where("vid=?", videoID).Update("comment_count", comCount)
	return result.Error
}

//插入新发布的视频
func InsertNewVideo(video *model.Video) error {
	result := middleware.GetMySQLClient().Create(video)
	return result.Error
}

//获取某用户发布的所有视频
func GetVideosByUserID(userID int64) []model.Video {
	var videos []model.Video
	middleware.GetMySQLClient().Where("uid=?", userID).Preload("User").Find(&videos)
	return videos
}

// 根据视频ID查询  点赞数量  和 评论数量
func GetVideoByID(videoID int64) model.Video {
	var video model.Video
	middleware.GetMySQLClient().Where("vid=?", videoID).Preload("User").Find(&video)
	return video
}
