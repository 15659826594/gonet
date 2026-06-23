package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var scopes = make(map[string]*viper.Viper)

func SetGlobalConfigFile(filename string) error {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		_ = viper.ReadInConfig()
	})
	return nil
}

func Viper() *Config {
	c := new(Config)
	err := viper.Unmarshal(c)
	if err != nil {
		log.Fatalf("Fatal unmarshal config file: %s\n", err)
	}
	return c
}

func Get(name string) (*viper.Viper, error) {
	if v, ok := scopes[name]; ok {
		return v, nil
	}
	return nil, errors.New("configuration file does not exist")
}

func SetConfigFile(name, filename string) (*viper.Viper, error) {
	if _, exist := scopes[name]; exist {
		return nil, errors.New("configuration already exists")
	}
	v := viper.New()
	v.SetEnvPrefix(name)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigFile(filename)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		_ = v.ReadInConfig()
	})
	scopes[name] = v
	return v, nil
}

type RenameFunc func(string) string

func LoadConfigGlob(pattern string, rename ...RenameFunc) {
	filenames, _ := doublestar.Glob(os.DirFS("./"), filepath.ToSlash(filepath.Clean(pattern)))

	renameFn := func(filename string) string {
		return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	}
	if len(rename) > 0 {
		renameFn = rename[0]
	}
	for _, filename := range filenames {
		_, _ = SetConfigFile(renameFn(filename), filename)
	}
}

func AllSettings() map[string]any {
	settings := viper.AllSettings()
	for s, v := range scopes {
		settings[s] = v.AllSettings()
	}
	return settings
}
