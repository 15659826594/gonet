package UserService

import (
	"horus/app/common/dao/User"
	"horus/app/common/model"
)

func GetById(id uint) *model.User {
	return UserDao.GetById(id)
}
