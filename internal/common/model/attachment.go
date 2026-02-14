package model

// Attachment 附件表
type Attachment struct {
	Id          uint   `gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Category    string `gorm:"column:category;size:50;comment:类别"`
	AdminId     uint   `gorm:"column:admin_id;not null;default:0;comment:管理员ID"`
	UserId      uint   `gorm:"column:user_id;not null;default:0;comment:会员ID"`
	Url         string `gorm:"column:url;size:255;comment:物理路径"`
	Imagewidth  uint   `gorm:"column:imagewidth;default:0;comment:宽度"`
	Imageheight uint   `gorm:"column:imageheight;default:0;comment:高度"`
	Imagetype   string `gorm:"column:imagetype;size:30;comment:图片类型"`
	Imageframes uint   `gorm:"column:imageframes;not null;default:0;comment:图片帧数"`
	Filename    string `gorm:"column:filename;size:100;comment:文件名称"`
	Filesize    uint   `gorm:"column:filesize;not null;comment:文件大小"`
	Mimetype    string `gorm:"column:mimetype;size:100;comment:mime类型"`
	Extparam    string `gorm:"column:extparam;size:255;comment:透传数据"`
	Createtime  int64  `gorm:"column:createtime;autoCreateTime;comment:创建日期"`
	Updatetime  int64  `gorm:"column:updatetime;autoUpdateTime;comment:更新时间"`
	Uploadtime  int64  `gorm:"column:uploadtime;comment:上传时间"`
	Storage     string `gorm:"column:storage;size:100;not null;default:local;comment:存储位置"`
	Sha1        string `gorm:"column:sha1;size:40;comment:文件 sha1编码"`
}
