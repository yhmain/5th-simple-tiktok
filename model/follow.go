package model

// 关注的操作
type Follow struct {
	Id       string `json:"id" gorm:"column:follow_id;primary_key"`      //关注主键
	UserA    User   `json:"ua" gorm:"ForeignKey:UserAID"`                //用户A
	UserAID  int64  `json:"ua_id" gorm:"column:uaid"`                    //外键：用户A的ID
	UserB    User   `json:"ub" gorm:"ForeignKey:UserBID"`                //用户B
	UserBID  int64  `json:"ub_id" gorm:"column:ubid"`                    //外键：用户B的ID
	IsFollow bool   `json:"is_follow" gorm:"column:is_follow;default:0"` //true:A关注B，false:A未关注B
}

func (f *Follow) TableName() string {
	return "tk_follow"
}
