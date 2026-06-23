package model

// Category 分类表
type Category struct {
	Id          uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Pid         uint   `gorm:"column:pid;not null;default:0;index;comment:父ID"`
	Type        string `gorm:"column:type;size:30;comment:栏目类型"`
	Name        string `gorm:"column:name;size:30"`
	Nickname    string `gorm:"column:nickname;size:50"`
	Flag        string `gorm:"column:flag;type:set('hot','index','recommend')"`
	Image       string `gorm:"column:image;size:100;comment:图片"`
	Keywords    string `gorm:"column:keywords;size:255;comment:关键字"`
	Description string `gorm:"column:description;size:255;comment:描述"`
	Diyname     string `gorm:"column:diyname;size:30;comment:自定义名称"`
	Createtime  int64  `gorm:"column:createtime;autoCreateTime;comment:创建时间"`
	Updatetime  int64  `gorm:"column:updatetime;autoUpdateTime;comment:更新时间"`
	Weigh       int    `gorm:"column:weigh;not null;default:0;index:,sort:desc;comment:权重"`
	Status      string `gorm:"column:status;size:30;comment:状态"`
}
