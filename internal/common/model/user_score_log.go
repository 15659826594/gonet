package model

// UserScoreLog 用户积分变动表
type UserScoreLog struct {
	Id         uint   `gorm:"column:id;primaryKey;autoIncrement"`
	UserId     uint   `gorm:"column:user_id;not null;default:0;comment:会员ID"`
	Score      int    `gorm:"column:score;not null;default:0;comment:变更积分"`
	Before     int    `gorm:"column:before;not null;default:0;comment:变更前积分"`
	After      int    `gorm:"column:after;not null;default:0;comment:变更后积分"`
	Memo       string `gorm:"column:memo;size:255;comment:备注"`
	Createtime int64  `gorm:"column:createtime;autoCreateTime;comment:创建时间"`
}
