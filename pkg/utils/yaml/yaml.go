package yaml

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Yaml struct {
	*yaml.Node
	filename string
}

func Load(filename string) (*Yaml, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	y := &Yaml{
		Node:     new(yaml.Node),
		filename: filename,
	}

	err = yaml.Unmarshal(bytes, y.Node)

	if err != nil {
		return nil, err
	}

	return y, nil
}

func (y *Yaml) Set(key string, value string) {
	keys := strings.Split(key, ".")
	// 查找并更新节点
	_ = set(y.Node, keys, value)
}

func set(node *yaml.Node, keys []string, value string) error {
	if len(keys) == 0 {
		return fmt.Errorf("keys不能为空")
	}

	// 如果是文档节点，进入其内容
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		return set(node.Content[0], keys, value)
	}

	// 如果是映射节点
	if node.Kind == yaml.MappingNode {
		targetKey := keys[0]

		// 遍历映射节点的内容
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			// 找到目标键
			if keyNode.Value == targetKey {
				// 如果是最后一个key，直接设置值
				if len(keys) == 1 {
					valueNode.SetString(value)
					return nil
				}
				// 否则继续递归查找下一级
				return set(valueNode, keys[1:], value)
			}
		}
		return fmt.Errorf("未找到键: %s", targetKey)
	}

	return fmt.Errorf("节点类型不支持更新: %v", node.Kind)
}

// removeBlankLines 移除YAML内容中的空白行
func removeBlankLines(data []byte) []byte {
	lines := bytes.Split(data, []byte("\n"))
	var filteredLines [][]byte

	for _, line := range lines {
		// 保留非空行（去除只包含空格的行）
		if len(bytes.TrimSpace(line)) > 0 {
			filteredLines = append(filteredLines, line)
		}
	}

	return bytes.Join(filteredLines, []byte("\n"))
}

func (y *Yaml) Save(filenames ...string) {
	filename := y.filename
	if len(filenames) > 0 {
		filename = filenames[0]
	}
	// 使用自定义编码器控制格式
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2) // 设置缩进
	err := encoder.Encode(y.Node)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	err = os.WriteFile(filename, removeBlankLines(buf.Bytes()), 0644)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
}
