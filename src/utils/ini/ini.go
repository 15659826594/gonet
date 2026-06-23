package ini

import (
	"fmt"
	"slices"

	"gopkg.in/ini.v1"
)

type Ini struct {
	*ini.File
	filename string
}

func Load(filename string) (*Ini, error) {
	file, err := ini.Load(filename)
	if err != nil {
		return nil, err
	}
	return &Ini{
		File:     file,
		filename: filename,
	}, nil
}

func (i *Ini) Set(key string, value string) {
	i.SetSection("", key, value)
}

func (i *Ini) SetMap(kv map[string]string) {
	i.SetSectionMap("", kv)
}

func (i *Ini) SetOrderMap(kv map[string]string, order []string) {
	i.SetSectionOrderMap("", kv, order)
}

func (i *Ini) SetSection(section, key string, value string) {
	i.File.Section(section).Key(key).SetValue(value)
}

func (i *Ini) SetSectionMap(section string, kv map[string]string) {
	sec := i.File.Section(section)
	for k, v := range kv {
		sec.Key(k).SetValue(v)
	}
}

func (i *Ini) SetSectionOrderMap(section string, kv map[string]string, order []string) {
	sec := i.File.Section(section)

	// 先按 order 顺序处理存在的键
	for _, k := range order {
		if v, ok := kv[k]; ok {
			sec.Key(k).SetValue(v)
		}
	}

	// 处理剩余未在 order 中指定的键
	if len(kv) > len(order) {
		for k, v := range kv {
			if !slices.Contains(order, k) {
				sec.Key(k).SetValue(v)
			}
		}
	}
}

func (i *Ini) Delete(keys ...string) {
	switch len(keys) {
	case 0:
		return
	case 1:
		i.Section("").DeleteKey(keys[0])
	default:
		section := i.Section(keys[0])
		section.DeleteKey(keys[1])
		if len(section.Keys()) == 0 && section.Name() != "" {
			i.File.DeleteSection(section.Name())
		}
	}
}

func (i *Ini) Save(filenames ...string) {
	filename := i.filename
	if len(filenames) > 0 {
		filename = filenames[0]
	}
	err := i.File.SaveTo(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
