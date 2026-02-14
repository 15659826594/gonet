package model

// UserToken 用户Token表
type UserToken struct {
	Token      string `gorm:"column:token;primaryKey;comment:Token"`
	UserId     uint   `gorm:"column:user_id;not null;default:0;comment:会员ID"`
	Createtime int64  `gorm:"column:createtime;autoCreateTime;comment:创建时间"`
	Expiretime int64  `gorm:"column:expiretime;comment:过期时间"`
}
