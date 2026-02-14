package model

// AdminLog 管理员日志表
type AdminLog struct {
	Id         uint   `gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	AdminId    uint   `gorm:"column:admin_id;not null;default:0;comment:管理员ID"`
	Username   string `gorm:"column:username;size:30;index;comment:管理员名字"`
	Url        string `gorm:"column:url;size:1500;comment:操作页面"`
	Title      string `gorm:"column:title;size:100;comment:日志标题"`
	Content    string `gorm:"column:content;type:longtext;not null;comment:内容"`
	Ip         string `gorm:"column:ip;size:50;comment:IP"`
	Useragent  string `gorm:"column:useragent;size:255;comment:User-Agent"`
	Createtime int64  `gorm:"column:createtime;autoCreateTime;comment:操作时间"`
}
