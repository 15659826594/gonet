package model

// Sms 短信验证码表
type Sms struct {
	Id         uint   `gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Event      string `gorm:"column:event;size:30;comment:事件"`
	Mobile     string `gorm:"column:mobile;size:20;comment:手机号"`
	Code       string `gorm:"column:code;size:10;comment:验证码"`
	Times      uint   `gorm:"column:times;not null;default:0;comment:验证次数"`
	Ip         string `gorm:"column:ip;size:30;comment:IP"`
	Createtime uint64 `gorm:"column:createtime;autoCreateTime;default:0;comment:创建时间"`
}
