package model

// AuthGroupAccess 权限组访问表
type AuthGroupAccess struct {
	Uid     uint `gorm:"column:uid;not null;comment:会员ID"`
	GroupId uint `gorm:"column:group_id;not null;comment:级别ID"`
}
