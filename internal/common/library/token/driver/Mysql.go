package driver

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"gota/internal/common/model"
	Config "gota/pkg/config"
	"gota/pkg/database"
	"hash"
	"time"

	"golang.org/x/crypto/ripemd160"
	"gorm.io/gorm"
)

type Mysql struct {
	handler *gorm.DB
	options map[string]any
}

func (t *Mysql) Construct(option map[string]any) *Mysql {
	t.options = map[string]any{
		"table":      "fa_user_token",
		"expire":     2592000,
		"connection": map[string]any{},
	}
	if option != nil {
		t.options = mapMerge(t.options, option)
	}
	if m, ok := t.options["connection"].(map[string]any); ok {
		if len(m) > 0 {

		} else {
			t.handler = database.Gorm().Table(t.options["table"].(string))
		}

		//times := time.Now().Unix()
		//t.handler.Where("expiretime < ? and expiretime > ?", times, 0)
	}
	return t
}

func (t *Mysql) Set(s string, i int, i2 int) {
	//TODO implement me
	panic("implement me")
}

func (t *Mysql) Get(token string) map[string]any {
	data := model.UserToken{}
	err := t.handler.Session(&gorm.Session{NewDB: true}).Where("token = ?", t.GetEncryptedToken(token)).Find(&data).Error
	if err == nil {
		if data.Expiretime != 0 && data.Expiretime > time.Now().Unix() {
			res := map[string]any{
				"token":      data.Token,
				"user_id":    data.UserId,
				"createtime": data.Createtime,
				"expiretime": data.Expiretime,
			}
			//返回未加密的token给客户端使用
			res["token"] = token
			//返回剩余有效时间
			res["expires_in"] = t.GetExpiredIn(data.Expiretime)
			return res
		} else {
			t.Delete(token)
		}
	}
	return nil
}

func (t *Mysql) Check(i int, i2 int) {
	//TODO implement me
	panic("implement me")
}

func (t *Mysql) Delete(token string) bool {
	//"2333" => "b54ea48a593f3f77b8f86bbf6c4c91d933a6e48c"
	t.handler.Session(&gorm.Session{NewDB: true}).Where("token = ?", t.GetEncryptedToken(token)).Delete(new(model.UserToken))
	return true
}

func (t *Mysql) Clear(s string) {
	//TODO implement me
	panic("implement me")
}

func (t *Mysql) Handler() *gorm.DB {
	return t.handler.Session(&gorm.Session{NewDB: true})
}

func (t *Mysql) GetEncryptedToken(token string) string {
	config := Config.Viper().Token
	return hashHmac(config.Hashalgo, token, config.Key)
}

func (t *Mysql) GetExpiredIn(expiretime int64) int64 {
	if expiretime == 0 {
		return 365 * 86400
	}
	return max(0, expiretime-time.Now().Unix())
}

func mapMerge[K comparable, V any](original, override map[K]V) map[K]V {
	result := make(map[K]V)

	// 先复制原始 map
	for k, v := range original {
		result[k] = v
	}

	// 再合并覆盖的 map，相同 key 会覆盖原始值
	for k, v := range override {
		result[k] = v
	}

	return result
}

func hashHmac(algo, data, key string) string {
	var m hash.Hash
	dataByte := []byte(data)
	keyByte := []byte(key)
	switch algo {
	case "sha256":
		m = hmac.New(sha256.New, keyByte)
	case "ripemd160":
		m = hmac.New(ripemd160.New, keyByte)
	default:
		panic("未定义加密类型")
	}
	_, _ = m.Write(dataByte)
	result := m.Sum(nil)
	return hex.EncodeToString(result)
}
