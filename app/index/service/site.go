package service

import (
	"horus/app/index/model/dto"
	"horus/src/database"
)

type Site struct{}

// GetSiteConfig 获取站点配置信息
func (s *Site) GetSiteConfig(names ...string) []dto.ConfigSiteDto {
	var sites []dto.ConfigSiteDto
	database.Gorm().Where("name IN ?", names).Find(&sites)
	return sites
}
