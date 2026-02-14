package vo

import (
	"gonet/internal/index/model/dto"
)

type ConfigSiteVo struct {
	Name  string `json:"name"`
	Beian string `json:"beian"`
}

func (r *ConfigSiteVo) Dto2Vo(dtos ...dto.ConfigSiteDto) ConfigSiteVo {
	for _, config := range dtos {
		switch config.Name {
		case "name":
			r.Name = config.Value
		case "beian":
			r.Beian = config.Value
		}
	}
	return *r
}
