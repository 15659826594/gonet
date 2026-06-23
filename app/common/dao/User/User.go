package UserDao

import (
	"horus/app/common/model"
	"horus/src/database"
)

func GetById(id uint) *model.User {
	user := &model.User{Id: id}
	result := database.Gorm().First(user)
	if result.Error == nil {
		return user
	}
	return nil
}
