package insp

import (
	"database/sql"
	"time"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/15
 */

type Result []map[string]interface{}

// RowsToMap 将 sql.Rows 转换为 []map[string]interface{} 结构
// 每行数据转为 map，键为列名，值为对应数据（自动处理类型转换）
// 返回结果切片及可能发生的错误
func RowsToMap(rows *sql.Rows) (Result, error) {
	// 获取列名集合
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 最终结果集
	var result []map[string]interface{}

	// 遍历每一行数据
	for rows.Next() {
		// 创建临时存储容器
		// 每个列值初始化为 interface{} 类型指针
		scanArgs := make([]interface{}, len(columns))
		for i := range scanArgs {
			scanArgs[i] = new(interface{})
		}

		// 扫描当前行数据
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		// 构建当前行映射
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			// 解引用指针获取实际值
			rawValue := *(scanArgs[i].(*interface{}))

			// 类型转换处理
			switch v := rawValue.(type) {
			case []byte:
				// 二进制数据转字符串
				rowMap[col] = string(v)
			case time.Time:
				// 时间类型标准化处理
				rowMap[col] = v.Format(time.RFC3339Nano)
			case nil:
				// 空值显式处理
				rowMap[col] = nil
			default:
				// 保留原始类型（int64/float64/string等）
				rowMap[col] = v
			}
		}

		result = append(result, rowMap)
	}

	// 检查迭代过程中的错误
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
