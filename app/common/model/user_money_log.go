package model

// UserMoneyLog 用户余额变动表
type UserMoneyLog struct {
	Id         uint    `gorm:"column:id;primaryKey;autoIncrement"`
	UserId     uint    `gorm:"column:user_id;not null;default:0;comment:会员ID"`
	Money      float64 `gorm:"column:money;type:decimal(10,2);not null;default:0.00;comment:变更余额"`
	Before     float64 `gorm:"column:before;type:decimal(10,2);not null;default:0.00;comment:变更前余额"`
	After      float64 `gorm:"column:after;type:decimal(10,2);not null;default:0.00;comment:变更后余额"`
	Memo       string  `gorm:"column:memo;size:255;comment:备注"`
	Createtime int64   `gorm:"column:createtime;autoCreateTime;comment:创建时间"`
}
