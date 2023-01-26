package model

//json其实对应前端的 键
//结构体名称+ID 即可设置外键
type Video struct {
	Id            int64  `json:"id" gorm:"column:vid;primary_key"`                  //视频主键
	Title         string `json:"title" gorm:"column:video_title"`                   //视频简介、标题
	PlayUrl       string `json:"play_url" gorm:"column:play_url"`                   //视频播放地址
	CoverUrl      string `json:"cover_url" gorm:"column:cover_url"`                 //视频封面地址
	FavoriteCount int64  `json:"favorite_count" gorm:"column:favorite_count"`       //视频点赞总数
	CommentCount  int64  `json:"comment_count" gorm:"column:comment_count"`         //视频评论总数
	CreatedTime   int64  `json:"created_time,omitempty" gorm:"column:created_time"` //视频创建时间，时间戳形式
	// IsFavorite    bool   `json:"is_favorite,omitempty" gorm:"column:is_favorite"`       //true:已点赞，false:未点赞
	User   User  `json:"author" gorm:"ForeignKey:UserID"` //视频作者
	UserID int64 `json:"user_id" gorm:"column:uid"`       //外键：视频作者的ID
}

//结构体Video对应数据库中videos表
func (v *Video) TableName() string {
	return "tk_video"
}
