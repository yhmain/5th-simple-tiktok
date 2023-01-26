package model

// omitempty: 如果字段值为空，则在编码期间忽略该字段
type User struct {
	Id            int64  `json:"id,omitempty" gorm:"column:uid;primary_key"`            //用户ID
	UserName      string `json:"name,omitempty" gorm:"uniqueIndex;column:user_name"`    //昵称（唯一索引）
	Password      string `json:"password,omitempty" gorm:"column:password"`             //密码
	FollowCount   int    `json:"follow_count,omitempty" gorm:"column:follow_count"`     //关注总数
	FollowerCount int    `json:"follower_count,omitempty" gorm:"column:follower_count"` //粉丝总数

	// IsFollow     bool `json:"is_follow,omitempty" gorm:"column:IsFollow"`           //true:已关注，false:未关注
	// CommentCount int  `json:"comment_count,omitempty" gorm:"column:CommentCount"`   //评论数目
	// IsFocused    bool `json:"is_favorite,omitempty" gorm:"column:IsFavorite"`       //是否喜欢
	// LikeCount    int  `json:"favorite_count,omitempty" gorm:"column:FavoriteCount"` //表示喜欢的视频数量
}

//结构体User对应数据库中tk_user表
func (u *User) TableName() string {
	return "tk_user"
}
