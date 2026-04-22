package UserService

import (
	"gota/internal/common/dao/User"
	"gota/internal/common/model"
)

func GetById(id uint) *model.User {
	return UserDao.GetById(id)
}
