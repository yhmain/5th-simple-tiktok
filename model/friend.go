package model

// 好友，此处假设id小的为A用户，大的为B用户
type Friend struct {
	Id       string `json:"id" gorm:"column:friend_id;primary_key"`      //好友主键
	UserA    User   `json:"ua" gorm:"ForeignKey:UserAID"`                //用户A
	UserAID  int64  `json:"ua_id" gorm:"column:uaid"`                    //外键：用户A的ID
	UserB    User   `json:"ub" gorm:"ForeignKey:UserBID"`                //用户B
	UserBID  int64  `json:"ub_id" gorm:"column:ubid"`                    //外键：用户B的ID
	IsFriend bool   `json:"is_friend" gorm:"column:is_friend;default:0"` //true:是，false:不是
}

func (f *Friend) TableName() string {
	return "tk_friend"
}
