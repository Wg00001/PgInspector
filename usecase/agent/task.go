package agent

import (
	"PgInspector/entities/agent"
	"PgInspector/entities/config"
	"PgInspector/entities/task"
	"PgInspector/usecase/agent/analyzer"
	"PgInspector/usecase/agent/format"
	"PgInspector/usecase/agent/kbase"
	"PgInspector/usecase/alerter"
	"PgInspector/usecase/logger"
	"fmt"
	"log"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/2/25
 */

type AgentTask struct {
	*config.AgentTaskConfig
}

func NewTask(taskConfig *config.AgentTaskConfig) *AgentTask {
	return &AgentTask{AgentTaskConfig: taskConfig}
}

var _ task.Task = (*AgentTask)(nil)

func (t *AgentTask) Do() error {
	//1. 获取日志
	contents, err := logger.Get(t.LogID).ReadLog(t.LogFilter)
	if err != nil {
		return err
	}

	//2. 组织日志格式(JSON)
	msg, err := format.Format(contents...)
	if err != nil {
		return err
	}
	if msg == nil {
		return fmt.Errorf("Ai task err\n- AgentTask name: %v\n- err: log read empty\n---\n", t.Name)
	}

	//3. 知识库检索 (并且组织格式)
	kbaseContent, err := t.KBaseSearch(msg)
	if err != nil {
		// 根据业务需求，不需要阻断流程
		log.Printf("Agent task warring: kbase search fail but continue to execute, Err :%v", err)
	}

	//4. 组织格式 (发生一次复制)
	context := &agent.AnalyzeContent{
		SystemMsg: t.SystemMessage,
		UserMsg:   *msg,
		KBaseMsg:  *kbaseContent,
	}

	//5. 发送Ai获取结果
	res, err := analyzer.Analyze(context)
	if err != nil {
		return err
	}

	//6. 将ai结果发送给Alert
	return alerter.GetAlert(t.AlertID).Send(*buildAiAlertContent(t, res))
}

func (t *AgentTask) GetCron() *config.Cron {
	return t.Cron
}

func (t *AgentTask) GetName() config.Name {
	return "ai_task:" + t.Name
}

func (t *AgentTask) KBaseSearch(msg *string) (*string, error) {
	if len(t.KBase) == 0 {
		return nil, fmt.Errorf("agent - kbase: agent has not kbase")
	}
	if msg == nil || *msg == "" {
		return nil, fmt.Errorf("empty input message")
	}
	//使用Ai生成日志的关键词
	query, err := generateQueryWithAI(msg)
	if err != nil {
		return nil, err
	}
	var kDocs []*agent.Document
	for _, kb := range t.KBase {
		kbaseObj := kbase.Get(kb)
		//根据KBase对应的嵌入向量生成器，生成相应的嵌入向量
		embeddingQuery, err := generateQueryEmbedding(query, kbaseObj)
		if err != nil {
			return nil, err
		}
		//进行关键词和向量的混合检索(自动去重)
		resDocs, err := hybridSearch(kbaseObj, query, embeddingQuery, t.KBaseTopN)
		if err != nil {
			return nil, err
		}
		kDocs = append(kDocs, resDocs...)
	}
	return formatKBaseContent(kDocs, t.KBaseMaxLen), nil
}
