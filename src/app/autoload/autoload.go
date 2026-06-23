package autoload

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

func init() {
	Load()
}

// Load 加载环境变量, 为.ini格式, 默认加载 .env
func Load(filenames ...string) []error {
	if len(filenames) == 0 {
		filenames = append(filenames, ".env")
	}
	var errs []error
	for _, filename := range filenames {
		env, err := ini.Load(filename)
		errs = append(errs, err)
		if err != nil {
			continue
		}
		for _, section := range env.Sections() {
			for _, key := range section.Keys() {
				os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", section.Name(), key.Name())), key.String())
			}
		}
	}
	return errs
}
