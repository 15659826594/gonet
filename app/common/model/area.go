package model

// Area 地区表
type Area struct {
	Id        uint   `gorm:"column:id;primaryKey;autoIncrement;comment:ID"`
	Pid       int    `gorm:"column:pid;comment:父id"`
	Shortname string `gorm:"column:shortname;size:100;comment:简称"`
	Name      string `gorm:"column:name;size:100;comment:名称"`
	Mergename string `gorm:"column:mergename;size:255;comment:全称"`
	Level     int8   `gorm:"column:level;comment:层级:1=省,2=市,3=区/县"`
	Pinyin    string `gorm:"column:pinyin;size:100;comment:拼音"`
	Code      string `gorm:"column:code;size:100;comment:长途区号"`
	Zip       string `gorm:"column:zip;size:100;comment:邮编"`
	First     string `gorm:"column:first;size:50;comment:首字母"`
	Lng       string `gorm:"column:lng;size:100;comment:经度"`
	Lat       string `gorm:"column:lat;size:100;comment:纬度"`
}
