package ai

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/22
 */

// Analyzer 无状态，接受Temporary并发送给ai，然后解析并返回结果
type Analyzer interface {
	Init() error
	// Analyze 将可供ai分析的string发送给AI进行分析，并解析返回结果。
	Analyze(string) (string, error)
}
