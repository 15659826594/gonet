package Auth

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

var (
	instance *Auth
	once     sync.Once
)

type Auth struct {
	*gin.Context
	*casbin.Enforcer
	error   error
	logined bool
	user    any
	token   string
	//Token默认有效时长
	keeptime   int
	requestUri string
	rules      []string
	//默认配置
	config      any
	options     any
	allowFields []string
}

func Instance() *Auth {
	once.Do(func() {
		_, filename, _, _ := runtime.Caller(0)
		dir := filepath.Dir(filename)
		var err error
		enforcer, err := casbin.NewEnforcer(filepath.Join(dir, "model.pml"), filepath.Join(dir, "policy.csv"))
		if err != nil {
			panic(err)
			return
		}
		instance = &Auth{
			Enforcer:    enforcer,
			allowFields: []string{"id", "username", "nickname", "mobile", "avatar", "score"},
		}
	})
	return instance
}

// GetUser 获取User模型
func (t *Auth) GetUser() any {
	return t.user
}

// Init 根据Token初始化
// 参数:
//
//	token: 令牌
//
// 返回值:
//
//	bool
func (t *Auth) Init(token string) bool {
	if t.logined {
		return true
	}
	if t.error != nil {
		return false
	}
	return true
}

// Register 注册用户
// 参数:
//
//	用户名: username
//	密码: password
//	邮箱: username
//	email: email
//	手机号: mobile
//	扩展参数: extend
//
// 返回值:
//
//	bool
func (t *Auth) Register(username, password, email, mobile string, extend map[string]any) bool {
	return true
}

// Login 用户登录
// 参数:
//
//	account: 账号,用户名、邮箱、手机号
//	password: 账号,密码
//
// 返回值:
//
//	bool
func (t *Auth) Login(account, password string) bool {
	return true
}

// Logout 退出
// 返回值:
//
//	bool
func (t *Auth) Logout(account, password string) bool {
	if !t.logined {
		t.SetError("You are not logged in")
		return false
	}
	//设置登录标识
	t.logined = false
	return true
}

// Changepwd 修改密码
// 参数:
//
//	newpassword: 新密码
//	oldpassword: 旧密码
//	ignoreoldpassword: 忽略旧密码
//
// 返回值:
//
//	bool
func (t *Auth) Changepwd(newpassword, oldpassword string, ignoreoldpassword bool) bool {
	if !t.logined {
		t.SetError("You are not logged in")
		return false
	}
	return true
}

// Direct 直接登录账号
// 参数:
//
//	user_id: user_id
//
// 返回值:
//
//	bool
func (t *Auth) Direct(userId string) {

}

// Check 检测是否是否有对应权限
// 参数:
//
//	path: 控制器/方法
//	module: 模块 默认为当前模块
//
// 返回值:
//
//	bool
func (t *Auth) Check(path, module string) bool {
	if !t.logined {
		return false
	}
	return true
}

// IsLogin 判断是否登录
// 返回值:
//
//	bool
func (t *Auth) IsLogin() bool {
	if t.logined {
		return true
	}
	return false
}

// GetToken 获取当前Token
// 返回值:
//
//	string
func (t *Auth) GetToken() string {
	return t.token
}

// GetUserinfo 获取会员基本信息
// 返回值:
//
//	Userinfo
func (t *Auth) GetUserinfo() any {
	return t.user
}

// GetRuleList 获取会员组别规则列表
// 返回值:
//
//	array|bool|\PDOStatement|string|\think\Collection
func (t *Auth) GetRuleList() any {
	return t.user
}

// GetRequestUri 获取当前请求的URI
// 返回值:
//
//	string
func (t *Auth) GetRequestUri() string {
	return t.requestUri
}

// SetRequestUri 设置当前请求的 URI
// 参数:
//
//	uri: 请求的 URI
func (t *Auth) SetRequestUri(uri string) {
	t.requestUri = uri
}

// GetAllowFields 获取允许输出的字段
// 返回值:
//
//	[]string
func (t *Auth) GetAllowFields() []string {
	return t.allowFields
}

// SetAllowFields 设置允许输出的字段
// 参数:
//
//	fields: []string
func (t *Auth) SetAllowFields(fields []string) {
	t.allowFields = fields
}

// Delete 删除一个指定会员
// 参数:
//
//	userId: 会员ID
// 返回值:
//
//	bool

func (t *Auth) Delete(userId int) bool {
	return true
}

func encrypMd5(s string) string {
	data := md5.Sum([]byte(s))
	return hex.EncodeToString(data[:])
}

// GetEncryptPassword 获取密码加密后的字符串
// 参数:
//
//	password: 密码
//	salt: 密码盐
// 返回值:
//
//	string

func (t *Auth) GetEncryptPassword(password, salt string) string {
	return encrypMd5(encrypMd5(password + salt))
}

// Match 检测当前控制器和方法是否匹配传递的数组
// 参数:
//
//	arr: 需要验证权限的数组
// 返回值:
//
//	bool

func (t *Auth) Match(arr []string) bool {
	return false
}

// Keeptime 设置会话有效时间
// 参数:
//
//	keeptime: 默认为永久

func (t *Auth) Keeptime(keeptime int) {
	t.keeptime = keeptime
}

// Render 渲染用户数据
// 参数:
//
//	datalist: 二维数组
//	fields: 加载的字段列表
//	fieldkey: 渲染的字段
//	renderkey: 结果字段

func (t *Auth) Render(datalist any, fields []string, fieldkey string, renderkey string) {
	if fieldkey == "" {
		fieldkey = "user_id"
	}
	if renderkey == "" {
		renderkey = "userinfo"
	}
}

// SetError 设置错误信息
// 参数:
//
//	err: 错误信息
func (t *Auth) SetError(err string) *Auth {
	t.error = errors.New(err)
	return t
}

// GetError 设置错误信息
// 返回值:
//
//	error
func (t *Auth) GetError() error {
	return t.error
}
