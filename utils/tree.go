package utils

import (
	"fmt"
	"strings"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/5
 */

type Tree map[string]interface{}

func (m Tree) Get(fullPath string) (string, error) {
	if m == nil {
		return "", fmt.Errorf("Tree is nil")
	}
	path := strings.Split(fullPath, ".")
	cur := m
	for i, p := range path {
		// 尝试从当前层级的映射中获取指定键的值
		value, exists := cur[p]
		if !exists {
			return "", fmt.Errorf("key %s not found at path level %d", p, i)
		}
		// 如果不是最后一个路径元素，且值是一个 Tree 类型，继续深入查找
		if i < len(path)-1 {
			nextMap, ok := value.(Tree)
			if !ok {
				return "", fmt.Errorf("expected a map at key %s, but got %T", p, value)
			}
			cur = nextMap
		} else {
			// 到达最后一个路径元素，尝试将值转换为字符串
			strValue, ok := value.(string)
			if !ok {
				return "", fmt.Errorf("expected a string at key %s, but got %T", p, value)
			}
			return strValue, nil
		}
	}
	return "", fmt.Errorf("unexpected error")
}
