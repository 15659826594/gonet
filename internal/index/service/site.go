package service

import (
	"gonet/internal/index/model/dto"
	"gonet/pkg/database"
)

type Site struct{}

// GetSiteConfig 获取站点配置信息
func (s *Site) GetSiteConfig(names ...string) []dto.ConfigSiteDto {
	var sites []dto.ConfigSiteDto
	database.Gorm().Where("name IN ?", names).Find(&sites)
	return sites
}
