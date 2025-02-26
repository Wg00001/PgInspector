package insp

import (
	"PgInspector/adapters/alerter_adapter"
	"PgInspector/entities/alerter"
	"PgInspector/entities/config"
	"PgInspector/entities/insp"
	alerter2 "PgInspector/usecase/alerter"
	"errors"
	"fmt"
	"strconv"
)

/**
 * @description: 用于构建insp的builder, 包括读取
 * @author Wg
 * @date 2025/1/19
 */

type NodeBuilder struct {
	insp.Node
	error
}

func (n NodeBuilder) Build() (insp.Node, error) {
	return n.Node, n.error
}

func (n NodeBuilder) WithName(name string) NodeBuilder {
	n.Name = name
	return n
}

func (n NodeBuilder) WithSQL(sql string) NodeBuilder {
	n.SQL = sql
	return n
}

func (n NodeBuilder) WithEmptyAlert() NodeBuilder {
	n.AlertFunc = func(content alerter.Content) error {
		return alerter_adapter.AlertEmpty{}.Send(content)
	}
	return n
}

func (n NodeBuilder) BuildAlertFunc(alertWhen string) NodeBuilder {
	n.AlertFunc, n.error = buildAlertFunc(alertWhen, n.AlertID)
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

// 告警函数生成器
func buildAlertFunc(alertWhen string, alertId config.ID) (func(alerter.Content) error, error) {
	condition, err := splitCondition(alertWhen)
	if err != nil {
		return nil, err
	}
	if len(condition) != 3 {
		return nil, fmt.Errorf("alerter err: invalid condition format")
	}

	field := condition[0]
	operator := condition[1]
	expectedValue := condition[2]

	return func(content alerter.Content) error {
		for _, row := range content.Result {
			// 获取实际值
			actualValue, exists := row[field]
			if !exists {
				return fmt.Errorf("alerter err: field %s not found", field)
			}

			// 类型转换处理
			comparisonResult, err := compareValues(actualValue, expectedValue, operator)
			if err != nil {
				return err
			}

			if comparisonResult { //如果触发报警条件，则发送报警
				fmt.Printf("[ALERT] Condition triggered: %s %s %s\n",
					field, operator, expectedValue)
				fmt.Printf("Current value: %v\n", actualValue)
				err = alerter2.GetAlert(alertId).Send(content.AddWhen(alertWhen)) //发送报警
				if err != nil {
					return err
				}
			}
		}
		return nil
	}, nil
}

// 通用值比较函数
func compareValues(actual interface{}, expected string, operator string) (bool, error) {
	// 统一转换为float64和string两种类型处理
	actualFloat, isFloat := tryConvertFloat(actual)
	expectedFloat, expectedIsFloat := tryConvertFloat(expected)

	// 数值比较
	if isFloat && expectedIsFloat {
		switch operator {
		case ">":
			return actualFloat > expectedFloat, nil
		case "<":
			return actualFloat < expectedFloat, nil
		case ">=":
			return actualFloat >= expectedFloat, nil
		case "<=":
			return actualFloat <= expectedFloat, nil
		case "==", "=":
			return actualFloat == expectedFloat, nil
		case "!=":
			return actualFloat != expectedFloat, nil
		default:
			return false, fmt.Errorf("unsupported operator: %s", operator)
		}
	}

	// 字符串比较
	actualStr := fmt.Sprintf("%v", actual)
	switch operator {
	case "==", "=":
		return actualStr == expected, nil
	case "!=":
		return actualStr != expected, nil
	case ">":
		return actualStr > expected, nil
	case "<":
		return actualStr < expected, nil
	default:
		return false, fmt.Errorf("operator %s not supported for string comparison", operator)
	}
}

// 尝试转换为数值类型
func tryConvertFloat(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err == nil {
			return f, true
		}
	}
	return 0, false
}
