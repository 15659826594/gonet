package common

import (
	"errors"
	"gonet/internal/common/model"
	"gonet/pkg/database"
	"gonet/pkg/utils"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Result struct {
	Code       int            `json:"code"`
	Msg        string         `json:"msg"`
	Time       int64          `json:"time"`
	Data       any            `json:"data"`
	Statuscode int            `json:"-"`
	Type       string         `json:"-"`
	Header     map[string]any `json:"-"`
}

func (r *Result) WithCode(code int) *Result {
	r.Code = code
	// 未设置状态码,根据code值判断
	if code >= 1000 || code < 200 {
		r.Statuscode = 200
	} else {
		r.Statuscode = code
	}
	return r
}

func (r *Result) WithType(_type string) *Result {
	t := _type
	if !slices.Contains([]string{"json", "xml", "jsonp"}, _type) {
		t = "json"
	}
	r.Type = t
	return r
}

func (r *Result) WithHeader(header map[string]any) *Result {
	r.Header = header
	return r
}

func Response(res *Result) {
	panic(res)
}
func Success(args ...any) *Result {
	res := &Result{
		Code:       1,
		Msg:        "",
		Time:       time.Now().Unix(),
		Data:       nil,
		Statuscode: http.StatusOK,
		Type:       "json",
	}
	if len(args) > 0 {
		if msg, ok := args[0].(string); ok {
			res.Msg = msg
		}
	}
	if len(args) > 1 {
		res.Data = args[1]
	}
	return res
}
func Error(args ...any) *Result {
	res := &Result{
		Code:       0,
		Msg:        "",
		Time:       time.Now().Unix(),
		Data:       nil,
		Statuscode: http.StatusOK,
		Type:       "json",
	}
	if len(args) > 0 {
		if msg, ok := args[0].(string); ok {
			res.Msg = msg
		}
	}
	if len(args) > 1 {
		res.Data = args[1]
	}
	return res
}

func Assign(c *gin.Context, key string, value any) {
	maps := c.GetStringMap("Think")
	if maps == nil {
		maps = map[string]any{}
	}
	maps[key] = value
	c.Set("Think", maps)
}

func GetUserByToken(t string) (*model.User, error) {
	if t == "" {
		return nil, errors.New("token不能为空")
	}
	t = GetEncryptedToken(t)
	token := new(model.UserToken)
	result := database.Gorm().Where("token = ? ", t).First(token)
	if result.Error != nil {
		return nil, errors.New("无效令牌")
	}
	if token.Expiretime < time.Now().Unix() {
		return nil, errors.New("令牌已过期")
	}
	user := new(model.User)
	result = database.Gorm().Where("id = ?", token.UserId).First(user)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}
	//异常用户
	if user.Status != "normal" {
		return nil, errors.New("用户状态异常")
	}
	return user, nil
}

// NewToken 生成一个新的token
func NewToken(token string, userId uint, expire int64) (string, error) {
	var expiretime int64 = 0
	if expire != 0 {
		expiretime = time.Now().Unix() + expire
	}
	encrypted := GetEncryptedToken(token)
	result := database.Gorm().Create(&model.UserToken{
		Token:      encrypted,
		UserId:     userId,
		Createtime: time.Now().Unix(),
		Expiretime: expiretime,
	})
	if result.Error != nil {
		return "", result.Error
	}
	return encrypted, nil
}

/*GetEncryptedToken
* 获取加密后的Token
* @param string $Token Token标识
* @return string
 */
func GetEncryptedToken(token string) string {
	hashalgo := viper.GetString("token.hashalgo")
	key := viper.GetString("token.key")
	return utils.HashHmac(hashalgo, token, key)
}
