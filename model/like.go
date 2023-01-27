package model

// 点赞视频的操作
type Like struct {
	Id         string `json:"id" gorm:"column:like_id;primary_key"`            //点赞主键
	User       User   `json:"author" gorm:"ForeignKey:UserID"`                 //用户
	UserID     int64  `json:"user_id" gorm:"column:uid"`                       //外键：用户的ID
	Video      Video  `json:"video" gorm:"ForeignKey:VideoID"`                 //视频
	VideoID    int64  `json:"video_id" gorm:"column:vid"`                      //外键：视频的ID
	IsFavorite bool   `json:"is_favorite" gorm:"column:is_favorite;default:0"` //true:已点赞，false:未点赞
}

func (u *Like) TableName() string {
	return "tk_like"
}
