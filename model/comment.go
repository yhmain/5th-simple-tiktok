package model

type Comment struct {
	Id          int64  `json:"id" gorm:"column:cid;primary_key"`       //评论id
	Content     string `json:"content" gorm:"column:content"`          //评论内容
	CreatedTime int64  `json:"create_date" gorm:"column:created_time"` //评论时间
	User        User   `json:"user" gorm:"ForeignKey:UserID"`          //创建该评论的用户id
	UserID      int64  `json:"user_id" gorm:"column:uid"`              //外键：发布评论的用户ID
	Video       Video  `json:"video" gorm:"ForeignKey:VideoID"`        //视频
	VideoID     int64  `json:"video_id" gorm:"column:vid"`             //外键：发布评论的视频ID
}

//结构体Comment对应数据库中comments表
func (c *Comment) TableName() string {
	return "tk_comment"
}
