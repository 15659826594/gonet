package token

import "gorm.io/gorm"

type Driver interface {
	Set(string, int, int)
	Get(string) map[string]any
	Check(int, int)
	Delete(string) bool
	Clear(string)
	Handler() *gorm.DB
	GetEncryptedToken(string) string
	GetExpiredIn(int64) int64
}
