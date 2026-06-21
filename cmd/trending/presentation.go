package main

import "strings"

// formatTopics 格式化仓库主题列表。
func formatTopics(topics []string) string {
	if len(topics) == 0 {
		return "未在 GitHub API 中公开 topics"
	}

	return strings.Join(topics, ", ")
}

// buildPurpose 生成仓库作用和解决问题说明。
func buildPurpose(repo repository) string {
	if repo.Description != "" {
		return repo.Description
	}
	if len(repo.ReadmeNotes) > 0 {
		return repo.ReadmeNotes[0]
	}
	return "该仓库进入 AI 相关 GitHub Trending Top 5，但公开简介较少，需要进入仓库 README 查看完整定位。"
}

// buildUseCases 生成仓库适用场景说明。
func buildUseCases(repo repository) []string {
	useCases := []string{
		"适合快速评估 GitHub AI 热榜中新出现或重新升温的技术方向，因为该仓库已获得短期社区关注。",
	}

	if hasTopic(repo, "mcp") || strings.Contains(strings.ToLower(repo.Description), "mcp") {
		useCases = append(useCases, "适合需要把外部工具、代码库、数据源接入 AI Agent 的场景，因为 MCP 能把能力封装成标准工具接口。")
	}
	if hasTopic(repo, "rag") || strings.Contains(strings.ToLower(repo.Description), "rag") {
		useCases = append(useCases, "适合知识库问答、文档检索和企业内部搜索场景，因为 RAG 能把私有数据补充进 LLM 上下文。")
	}
	if hasTopic(repo, "agent") || hasTopic(repo, "agents") || strings.Contains(strings.ToLower(repo.Description), "agent") {
		useCases = append(useCases, "适合多步骤自动化、工具调用和复杂任务编排场景，因为 Agent 模式能把规划、执行、观察和修正串起来。")
	}
	if hasTopic(repo, "skill") || hasTopic(repo, "skills") || strings.Contains(strings.ToLower(repo.Description), "skill") {
		useCases = append(useCases, "适合团队沉淀可复用 AI 能力的场景，因为 Skill 把提示词、工具和流程封装成可发现、可组合的单元。")
	}
	if len(useCases) == 1 {
		useCases = append(useCases, "适合围绕 "+formatTopics(repo.Topics)+" 做技术调研、竞品分析或原型验证，因为仓库主题与当前 "+activeTopicName()+" 热点高度相关。")
	}

	return useCases
}

// buildArchitectureThoughts 生成仓库架构思想说明。
func buildArchitectureThoughts(repo repository) []string {
	thoughts := []string{
		"它成为热榜的核心原因通常不是单点功能，而是把模型能力、工具、数据和工作流组织成更容易落地的工程结构。",
	}
	if repo.Stars != "" {
		thoughts = append(thoughts, "当前 Stars 为 "+repo.Stars+"，说明它不只是概念验证，还积累了可观的社区验证和传播势能。")
	}
	if len(repo.Topics) > 0 {
		thoughts = append(thoughts, "相比只提供单一脚本的仓库，它用 "+strings.Join(repo.Topics, ", ")+" 等 topics 明确了能力边界，更容易被目标用户检索和采用。")
	}
	if repo.Language != "" {
		thoughts = append(thoughts, "使用 "+repo.Language+" 作为主要实现语言，降低了对应生态开发者集成、扩展和二次开发的成本。")
	}
	thoughts = append(thoughts, "它的稀缺性在于把热门 "+activeTopicName()+" 能力包装成可运行、可组合、可观察的工程入口，而不是停留在论文、提示词或孤立 Demo。")

	return thoughts
}

// hasTopic 判断仓库是否包含指定 topic。
func hasTopic(repo repository, topic string) bool {
	for _, repoTopic := range repo.Topics {
		if strings.EqualFold(repoTopic, topic) {
			return true
		}
	}

	return false
}

// buildPrinciples 生成仓库原理和实现思路说明。
func buildPrinciples(repo repository) []string {
	principles := make([]string, 0, 4)
	if len(repo.ReadmeNotes) > 0 {
		principles = append(principles, repo.ReadmeNotes...)
	}
	if len(principles) == 0 {
		principles = append(principles, "未能从 README 中自动抽取到足够的原理说明，建议进入仓库查看 architecture、how it works 或 design 文档。")
	}
	principles = append(principles, "以上内容由 GitHub 公开 README 自动摘取和归纳，适合作为快速了解入口，深入实现仍以仓库源码和文档为准。")

	return principles
}

// buildArchitectureDiagram 生成仓库原理和实现思路的 Mermaid 架构图。
func buildArchitectureDiagram(repo repository) string {
	switch {
	case hasTopic(repo, "mcp") || strings.Contains(strings.ToLower(repo.Description), "mcp"):
		return strings.Join([]string{
			"flowchart LR",
			"    User[用户 / AI 编程助手] --> Client[Agent Client]",
			"    Client --> Protocol[MCP 协议层]",
			"    Protocol --> Server[" + mermaidLabel(repo.Name) + "]",
			"    Server --> Tools[工具接口 / Skills]",
			"    Server --> Index[代码索引 / 知识图谱]",
			"    Server --> Data[文件系统 / API / 数据源]",
			"    Tools --> Result[结构化结果]",
			"    Index --> Result",
			"    Data --> Result",
			"    Result --> Client",
			"    Client --> Answer[生成回答 / 执行动作]",
		}, "\n")
	case hasTopic(repo, "rag") || strings.Contains(strings.ToLower(repo.Description), "rag"):
		return strings.Join([]string{
			"flowchart LR",
			"    User[用户问题] --> Query[查询理解]",
			"    Query --> Retriever[检索器]",
			"    Docs[文档 / 知识库] --> Chunk[切分与向量化]",
			"    Chunk --> Store[向量库 / 索引]",
			"    Store --> Retriever",
			"    Retriever --> Context[相关上下文]",
			"    Context --> LLM[LLM 生成器]",
			"    Query --> LLM",
			"    LLM --> Eval[引用校验 / 质量评估]",
			"    Eval --> Answer[可信答案]",
		}, "\n")
	case hasTopic(repo, "skill") || hasTopic(repo, "skills") || strings.Contains(strings.ToLower(repo.Description), "skill"):
		return strings.Join([]string{
			"flowchart LR",
			"    User[用户任务] --> Planner[Agent 任务规划]",
			"    Planner --> Registry[Skill 注册表]",
			"    Registry --> Select[能力匹配 / 权限校验]",
			"    Select --> Skill[可复用 Skill]",
			"    Skill --> Tool[工具 / API / Prompt]",
			"    Tool --> Observation[执行结果]",
			"    Observation --> Planner",
			"    Planner --> Output[最终交付]",
		}, "\n")
	case hasTopic(repo, "agent") || hasTopic(repo, "agents") || strings.Contains(strings.ToLower(repo.Description), "agent"):
		return strings.Join([]string{
			"flowchart LR",
			"    User[用户目标] --> Planner[任务规划 Agent]",
			"    Planner --> Memory[上下文记忆]",
			"    Planner --> Tools[工具调用层]",
			"    Tools --> APIs[外部 API / 本地工具]",
			"    APIs --> Observation[观察结果]",
			"    Observation --> Critic[反思 / 评估]",
			"    Critic --> Planner",
			"    Planner --> Deliverable[最终结果]",
		}, "\n")
	default:
		return strings.Join([]string{
			"flowchart LR",
			"    User[用户需求] --> Interface[应用入口]",
			"    Interface --> Orchestrator[AI 编排层]",
			"    Orchestrator --> Model[LLM / 模型能力]",
			"    Orchestrator --> Data[领域数据 / 上下文]",
			"    Orchestrator --> Tools[工具与自动化流程]",
			"    Model --> Result[候选结果]",
			"    Data --> Result",
			"    Tools --> Result",
			"    Result --> Review[校验 / 观测 / 反馈]",
			"    Review --> Output[可交付结果]",
		}, "\n")
	}
}

// mermaidLabel 清理 Mermaid 节点标签中的特殊字符。
func mermaidLabel(label string) string {
	label = strings.ReplaceAll(label, "[", "(")
	label = strings.ReplaceAll(label, "]", ")")
	label = strings.ReplaceAll(label, "\n", " ")
	return label
}

// emptyToDash 将空字符串转换为占位符。
func emptyToDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}
