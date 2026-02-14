package model

// UserGroup 用户组表
type UserGroup struct {
	Id         uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name       string `gorm:"column:name;size:50;comment:组名"`
	Rules      string `gorm:"column:rules;comment:权限节点"`
	Createtime int64  `gorm:"column:createtime;autoCreateTime;comment:添加时间"`
	Updatetime int64  `gorm:"column:updatetime;autoUpdateTime;comment:更新时间"`
	Status     string `gorm:"column:status;type:enum('normal','hidden');comment:状态"`
}
