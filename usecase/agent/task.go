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
	"context"
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

func (t *AgentTask) Do(context.Context) error {
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
	//3. Ai生成关键词 + 知识库检索 (并且组织格式)
	kbaseContent, err := t.KBaseSearch(msg)
	if err != nil || kbaseContent == nil {
		// 根据业务需求，不需要阻断流程
		log.Println("Agent task warring: kbase search fail but continue to execute, Err :%v", err)
		f := "知识库无相关内容\n"
		kbaseContent = &f
	}
	//4. 组织格式 (发生一次复制)
	content := &agent.AnalyzeContent{
		SystemMsg: t.SystemMessage,
		UserMsg:   *msg,
		KBaseMsg:  *kbaseContent,
	}

	//5. 发送Ai获取结果
	res, err := analyzer.Analyze(content)
	if err != nil {
		return err
	}

	//6.1 自学习（将巡检结果发进知识库）

	//6.2 将ai结果发送给Alert
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
	queryData, err := generateQueryWithAI(msg)
	if err != nil {
		return nil, err
	}

	var kDocs []*agent.Document
	for _, kb := range t.KBase {
		kbaseObj := kbase.Get(kb)
		//根据KBase对应的嵌入向量生成器，生成相应的嵌入向量
		//embeddingQuery, err := kbaseObj.Embedding(query)
		//if err != nil {
		//	return nil, err
		//}
		if kbaseObj == nil {
			return nil, fmt.Errorf("agent task : kbase not exist")
		}
		//进行关键词和向量的混合检索(自动去重)
		resDocs, err := kbaseObj.Search(*queryData)
		//resDocs, err := hybridSearch(kbaseObj, query, embeddingQuery, t.KBaseResults)
		if err != nil {
			return nil, err
		}
		kDocs = append(kDocs, resDocs...)
	}
	return formatKBaseContent(kDocs, t.KBaseMaxLen), nil
}

func (t *AgentTask) K() {
	//置信度评估
	//关键词提取
	//对比去重
	//人工审核
	//入库
}
