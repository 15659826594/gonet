package model

// AuthGroup 权限组表
type AuthGroup struct {
	Id         uint   `gorm:"column:id;primaryKey;autoIncrement"`
	PId        uint   `gorm:"column:pid;not null;default:0;comment:父组别"`
	Name       string `gorm:"column:name;size:100;comment:组名"`
	Rules      string `gorm:"column:rules;type:text;not null;comment:规则ID"`
	Createtime int64  `gorm:"column:createtime;autoCreateTime;comment:创建时间"`
	Updatetime int64  `gorm:"column:updatetime;autoUpdateTime;comment:更新时间"`
	Status     string `gorm:"column:status;size:30;comment:状态"`
}
