package model

// Version 版本表
type Version struct {
	Id          uint   `gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Oldversion  string `gorm:"column:oldversion;size:30;comment:旧版本号"`
	Newversion  string `gorm:"column:newversion;size:30;comment:新版本号"`
	Packagesize string `gorm:"column:packagesize;size:30;comment:包大小"`
	Content     string `gorm:"column:content;size:500;comment:升级内容"`
	Downloadurl string `gorm:"column:downloadurl;size:255;comment:下载地址"`
	Enforce     uint8  `gorm:"column:enforce;not null;default:0;comment:强制更新"`
	Createtime  int64  `gorm:"column:createtime;autoCreateTime;comment:创建时间"`
	Updatetime  int64  `gorm:"column:updatetime;autoUpdateTime;comment:更新时间"`
	Weigh       int    `gorm:"column:weigh;default:0;comment:权重"`
	Status      string `gorm:"column:status;size:30;comment:状态"`
}
