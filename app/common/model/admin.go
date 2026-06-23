package model

// Admin 管理员表
type Admin struct {
	Id           uint   `gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Username     string `gorm:"column:username;size:20;uniqueIndex;comment:用户名"`
	Nickname     string `gorm:"column:nickname;size:50;comment:昵称"`
	Password     string `gorm:"column:password;size:32;comment:密码"`
	Salt         string `gorm:"column:salt;size:30;comment:密码盐"`
	Avatar       string `gorm:"column:avatar;size:255;comment:头像"`
	Email        string `gorm:"column:email;size:100;comment:电子邮箱"`
	Mobile       string `gorm:"column:mobile;size:11;comment:手机号码"`
	Loginfailure uint8  `gorm:"column:loginfailure;not null;default:0;comment:失败次数"`
	Logintime    int64  `gorm:"column:logintime;comment:登录时间"`
	Loginip      string `gorm:"column:loginip;size:50;comment:登录IP"`
	Createtime   int64  `gorm:"column:createtime;autoCreateTime;comment:创建时间"`
	Updatetime   int64  `gorm:"column:updatetime;autoUpdateTime;comment:更新时间"`
	Token        string `gorm:"column:token;size:59;comment:Session标识"`
	Status       string `gorm:"column:status;size:30;not null;default:normal;comment:状态"`
}
