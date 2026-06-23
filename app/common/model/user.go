package model

// User 用户表
type User struct {
	Id               uint    `gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	GroupId          uint    `gorm:"column:group_id;not null;default:0;comment:组别ID"`
	Username         string  `gorm:"column:username;size:32;index;comment:用户名"`
	Nickname         string  `gorm:"column:nickname;size:50;comment:昵称"`
	Password         string  `gorm:"column:password;size:32;comment:密码"`
	Salt             string  `gorm:"column:salt;size:30;comment:密码盐"`
	Email            string  `gorm:"column:email;size:100;index;comment:电子邮箱"`
	Mobile           string  `gorm:"column:mobile;size:11;index;comment:手机号"`
	Avatar           string  `gorm:"column:avatar;size:255;comment:头像"`
	Level            uint8   `gorm:"column:level;not null;default:0;comment:等级"`
	Gender           uint8   `gorm:"column:gender;not null;default:0;comment:性别"`
	Birthday         string  `gorm:"column:birthday;type:date;comment:生日"`
	Bio              string  `gorm:"column:bio;size:100;comment:格言"`
	Money            float64 `gorm:"column:money;type:decimal(10,2);not null;default:0.00;comment:余额"`
	Score            int     `gorm:"column:score;not null;default:0;comment:积分"`
	Successions      uint    `gorm:"column:successions;not null;default:1;comment:连续登录天数"`
	Maxsuccessions   uint    `gorm:"column:maxsuccessions;not null;default:1;comment:最大连续登录天数"`
	Prevtime         int64   `gorm:"column:prevtime;comment:上次登录时间"`
	Logintime        int64   `gorm:"column:logintime;comment:登录时间"`
	Loginip          string  `gorm:"column:loginip;size:50;comment:登录IP"`
	Loginfailure     uint8   `gorm:"column:loginfailure;not null;default:0;comment:失败次数"`
	Loginfailuretime int64   `gorm:"column:loginfailuretime;comment:最后登录失败时间"`
	Joinip           string  `gorm:"column:joinip;size:50;comment:加入IP"`
	Jointime         int64   `gorm:"column:jointime;comment:加入时间"`
	Createtime       int64   `gorm:"column:createtime;autoCreateTime;comment:创建时间"`
	Updatetime       int64   `gorm:"column:updatetime;autoUpdateTime;comment:更新时间"`
	Token            string  `gorm:"column:token;size:50;comment:Token"`
	Status           string  `gorm:"column:status;size:30;comment:状态"`
	Verification     string  `gorm:"column:verification;size:255;comment:验证"`
}
