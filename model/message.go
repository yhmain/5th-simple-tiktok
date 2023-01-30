package model

// 消息
type Message struct {
	Id          int64  `json:"id" gorm:"column:message_id;primary_key"` //消息主键
	Content     string `json:"content" gorm:"column:content"`           //消息内容
	CreatedTime string `json:"create_time" gorm:"column:created_time"`  //消息时间，格式： yyyy-MM-dd HH:MM:ss
	UserA       User   `json:"ua" gorm:"ForeignKey:UserAID"`            //用户A
	UserAID     int64  `json:"ua_id" gorm:"column:uaid"`                //外键：用户A的ID
	UserB       User   `json:"ub" gorm:"ForeignKey:UserBID"`            //用户B
	UserBID     int64  `json:"ub_id" gorm:"column:ubid"`                //外键：用户B的ID
}

func (f *Message) TableName() string {
	return "tk_message"
}
