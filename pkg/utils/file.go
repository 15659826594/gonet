package utils

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

// FileListener 监听文件变动, callback返回true时取消监听
func FileListener(dir string, callback func(fsnotify.Event, func())) error {
	// 创建文件监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("创建文件监听器失败: %v", err)
	}

	// 添加监听目录
	err = watcher.Add(dir)
	if err != nil {
		watcher.Close()
		return fmt.Errorf("添加监听目录失败: %v", err)
	}

	// 控制监听器结束的标志
	doneChan := make(chan struct{})

	// 启动goroutine监听事件
	go func() {
		defer watcher.Close()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// 创建done函数，用于结束监听
				done := func() {
					close(doneChan)
				}

				if callback != nil {
					callback(event, done)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("监听器错误: %v", err)
				return
			case <-doneChan:
				return
			}
		}
	}()

	return nil
}

//// UpdateYamlNode 通过yaml.Node 更新 Yaml 文件
//func UpdateYamlNode(filename, key string, value string) error {
//	siteBytes, err := os.ReadFile(filename)
//
//	if err != nil {
//		return err
//	}
//
//	var node yaml.Node
//
//	err = yaml.Unmarshal(siteBytes, &node)
//
//	if err != nil {
//		return err
//	}
//
//	// 按"."分割key路径
//	keys := strings.Split(key, ".")
//
//	// 查找并更新节点
//	err = updateYamlNode(&node, keys, value)
//	if err != nil {
//		return err
//	}
//	// 使用自定义编码器控制格式
//	var buf bytes.Buffer
//	encoder := yaml.NewEncoder(&buf)
//	encoder.SetIndent(2) // 设置缩进
//	err = encoder.Encode(&node)
//	if err != nil {
//		return err
//	}
//
//	// 写回文件
//	err = os.WriteFile(filename, buf.Bytes(), 0644)
//	return err
//}
//
//// updateYamlNode 递归查找并更新YAML节点
//func updateYamlNode(node *yaml.Node, keys []string, value string) error {
//	if len(keys) == 0 {
//		return fmt.Errorf("keys不能为空")
//	}
//
//	// 如果是文档节点，进入其内容
//	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
//		return updateYamlNode(node.Content[0], keys, value)
//	}
//
//	// 如果是映射节点
//	if node.Kind == yaml.MappingNode {
//		targetKey := keys[0]
//
//		// 遍历映射节点的内容
//		for i := 0; i < len(node.Content); i += 2 {
//			keyNode := node.Content[i]
//			valueNode := node.Content[i+1]
//
//			// 找到目标键
//			if keyNode.Value == targetKey {
//				// 如果是最后一个key，直接设置值
//				if len(keys) == 1 {
//					valueNode.SetString(value)
//					return nil
//				}
//				// 否则继续递归查找下一级
//				return updateYamlNode(valueNode, keys[1:], value)
//			}
//		}
//		return fmt.Errorf("未找到键: %s", targetKey)
//	}
//
//	return fmt.Errorf("节点类型不支持更新: %v", node.Kind)
//}
