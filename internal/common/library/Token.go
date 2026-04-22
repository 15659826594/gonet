package library

import (
	"encoding/json"
	"gota/internal/common/library/token"
	"gota/internal/common/library/token/driver"
	Config "gota/pkg/config"
)

var (
	instance map[string]token.Driver
	handler  token.Driver
)

type Token struct {
}

func (t *Token) connect(options ...map[string]any) token.Driver {
	option := map[string]any{}
	if len(options) > 0 {
		option = options[0]
	}

	switch option["Type"] {
	case "Mysql":
		handler = new(driver.Mysql).Construct(option)
	case "Redis":
		//handler = new(driver.Redis).Construct(option)
	}
	return handler
}

func (t *Token) init(options ...map[string]any) token.Driver {
	var option map[string]any
	if len(options) > 0 {
		option = options[0]
	}
	if handler == nil {
		if option == nil && Config.Viper().Token.Type == "complex" {
			//$default = Config::get('token.default');
			//// 获取默认Token配置，并连接
			//$options = Config::get('token.' . $default['type']) ?: $default;
		} else if option == nil {
			option, _ = structToMap(Config.Viper().Token)
		}

		handler = t.connect(option)
	}
	return handler
}

func (t *Token) get(token string) any {
	return nil
}
func Get(token string) map[string]any {
	return new(Token).init().Get(token)
}

func DefaultGet(token string, def any) any {
	ret := new(Token).init().Get(token)
	if ret != nil {
		return ret
	}
	return def
}

func structToMap(s any) (map[string]any, error) {
	var result map[string]any

	// 先序列化为 JSON
	jsonData, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	// 再反序列化为 map
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
