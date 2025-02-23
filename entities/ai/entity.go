package ai

import "PgInspector/entities/db"

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/22
 */

// Analyzer 无状态，接受Temporary并发送给ai，然后解析并返回结果
type Analyzer interface {
	Init() error
	Analyze(Temporary) (string, error) //将暂时保存的巡检结果发给ai分析
}

// Temporary 无状态，将巡检结果进行暂存和读取，不应该过大
type Temporary interface {
	Init() error
	Append(db.Result) error //将巡检结果暂时保存
	Read() (string, error)  //获取暂存结果
	Clear() error           //清空暂存文件
}
