package model

// AuthRule 权限规则表
type AuthRule struct {
	Id         uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Type       string `gorm:"column:type;type:enum('menu','file');not null;default:file;comment:menu为菜单,file为权限节点"`
	PId        uint   `gorm:"column:pid;not null;default:0;index;comment:父ID"`
	Name       string `gorm:"column:name;size:100;uniqueIndex;comment:规则名称"`
	Title      string `gorm:"column:title;size:50;comment:规则名称"`
	Icon       string `gorm:"column:icon;size:50;comment:图标"`
	Url        string `gorm:"column:url;size:255;comment:规则URL"`
	Condition  string `gorm:"column:condition;size:255;comment:条件"`
	Remark     string `gorm:"column:remark;size:255;comment:备注"`
	Ismenu     uint8  `gorm:"column:ismenu;not null;default:0;comment:是否为菜单"`
	Menutype   string `gorm:"column:menutype;type:enum('addtabs','blank','dialog','ajax');comment:菜单类型"`
	Extend     string `gorm:"column:extend;size:255;comment:扩展属性"`
	Py         string `gorm:"column:py;size:30;comment:拼音首字母"`
	Pinyin     string `gorm:"column:pinyin;size:100;comment:拼音"`
	Createtime int64  `gorm:"column:createtime;autoCreateTime;comment:创建时间"`
	Updatetime int64  `gorm:"column:updatetime;autoUpdateTime;comment:更新时间"`
	Weigh      int    `gorm:"column:weigh;default:0;index;comment:权重"`
	Status     string `gorm:"column:status;size:30;comment:状态"`
}
