package model

// UserRule 用户规则表
type UserRule struct {
	Id         uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Pid        int    `gorm:"column:pid;comment:父ID"`
	Name       string `gorm:"column:name;size:50;comment:名称"`
	Title      string `gorm:"column:title;size:50;comment:标题"`
	Remark     string `gorm:"column:remark;size:100;comment:备注"`
	Ismenu     bool   `gorm:"column:ismenu;comment:是否菜单"`
	Createtime int64  `gorm:"column:createtime;autoCreateTime;comment:创建时间"`
	Updatetime int64  `gorm:"column:updatetime;autoUpdateTime;comment:更新时间"`
	Weigh      int    `gorm:"column:weigh;default:0;comment:权重"`
	Status     string `gorm:"column:status;type:enum('normal','hidden');comment:状态"`
}
