package model

// Ems 邮箱验证码表
type Ems struct {
	Id         uint   `gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Event      string `gorm:"column:event;size:30;comment:事件"`
	Email      string `gorm:"column:email;size:100;comment:邮箱"`
	Code       string `gorm:"column:code;size:10;comment:验证码"`
	Times      uint   `gorm:"column:times;not null;default:0;comment:验证次数"`
	Ip         string `gorm:"column:ip;size:30;comment:IP"`
	Createtime int64  `gorm:"column:createtime;autoCreateTime;comment:创建时间"`
}
