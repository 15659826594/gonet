package dto

import (
	"gonet/pkg/database"
)

type ConfigSiteDto struct {
	Name  string `gorm:"column:name;size:30;uniqueIndex;comment:变量名"`
	Value string `gorm:"column:value;comment:变量值"`
}

func (c *ConfigSiteDto) TableName() string {
	return database.GetTablePrefix() + "config"
}
