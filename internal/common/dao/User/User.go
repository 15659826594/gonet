package UserDao

import (
	"gota/internal/common/model"
	"gota/pkg/database"
)

func GetById(id uint) *model.User {
	user := &model.User{Id: id}
	result := database.Gorm().First(user)
	if result.Error == nil {
		return user
	}
	return nil
}
