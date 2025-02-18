package insp

import (
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
	"errors"
	"fmt"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/1/19
 */

const (
	keyAlertId   = "_alertid"
	keyAlertWhen = "_alertwhen"
)

type NodeBuilder struct {
	insp.Node
	error
}

func (n NodeBuilder) WithName(name string) NodeBuilder {
	n.Name = name
	return n
}

func (n NodeBuilder) WithSQL(sql string) NodeBuilder {
	n.SQL = sql
	return n
}

func (n NodeBuilder) ParseMap(arg map[string]interface{}) NodeBuilder {
	defer func() {
		if r := recover(); r != nil {
			n.error = fmt.Errorf("inspect node build fail, please check inspect config \nerr: %v\n", r)
		}
	}()
	alertId, ok := arg[keyAlertId]
	if !ok {
		return n
	} else {
		n.AlertID = config.ID(alertId.(int))
		delete(arg, keyAlertId)
	}
	alertWhen, ok := arg[keyAlertWhen]
	if !ok {
		return n
	} else {
		delete(arg, keyAlertWhen)
	}
	condition, err := splitCondition(alertWhen.(string))
	if err != nil {
		n.error = err
		return n
	}
	_ = condition
	n.AlertFunc = func(result insp.Result) error {
		return nil
	}
	return n
}

func splitCondition(s string) ([]string, error) {
	// 按优先级排序的运算符列表（长运算符优先）
	operators := []string{">=", "<=", "!=", "==", ">", "<", "="}

	// 遍历每个字符位置
	for i := 0; i < len(s); i++ {
		// 检查所有可能的运算符
		for _, op := range operators {
			opLen := len(op)
			// 检查是否超出字符串长度
			if i+opLen > len(s) {
				continue
			}
			// 匹配运算符
			if s[i:i+opLen] == op {
				left := s[:i]
				right := s[i+opLen:]
				// 验证左右操作数非空
				if left == "" || right == "" {
					return nil, errors.New("invalid format: missing operands")
				}
				return []string{left, op, right}, nil
			}
		}
	}
	return nil, errors.New("no valid operator found")
}
