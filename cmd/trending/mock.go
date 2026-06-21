package main

// mockRepositories 返回用于本地预览页面效果的固定仓库数据。
func mockRepositories() []repository {
	return []repository{
		{
			Name:        "microsoft/autogen",
			URL:         "https://github.com/microsoft/autogen",
			Description: "面向内容团队的 AI 视频 Agent 工作台，用自然语言完成脚本、素材生成、剪辑编排和导出。",
			Language:    "TypeScript",
			Stars:       "18,420",
			Topics:      []string{"ai", "agent", "video-generation", "workflow", "llm", "typescript"},
			ReadmeNotes: []string{
				"系统把视频生产拆成脚本规划、镜头设计、素材生成、时间线编排和导出五个阶段，每个阶段由独立 Agent 负责。",
				"核心原理是用 LLM 生成结构化制作计划，再通过工具调用连接图像生成、语音合成、字幕和剪辑引擎。",
				"适合验证 AI 内容生产类项目在 README 和网页中展示长描述、topics 和 Star 趋势图的效果。",
			},
		},
		{
			Name:        "modelcontextprotocol/servers",
			URL:         "https://github.com/modelcontextprotocol/servers",
			Description: "为 AI 编程助手提供代码库长期记忆的 MCP Server，解决大仓库上下文不足和重复索引问题。",
			Language:    "Go",
			Stars:       "32,105",
			Topics:      []string{"mcp", "code-intelligence", "knowledge-graph", "agent", "go"},
			ReadmeNotes: []string{
				"项目会扫描代码符号、调用关系和文档片段，将结构化信息写入本地知识图谱，供 AI Agent 按需查询。",
				"关键设计是把一次性上下文输入改成可复用的增量索引，减少 token 消耗并提升跨文件问题定位能力。",
				"运行时通过 MCP 协议暴露 search、explain、dependency 等工具，让 Claude、Cursor 等客户端直接调用。",
			},
		},
		{
			Name:        "langchain-ai/langchain",
			URL:         "https://github.com/langchain-ai/langchain",
			Description: "在日志、RAG 片段和工具输出进入 LLM 前进行语义压缩，降低 token 成本并保留关键事实。",
			Language:    "Python",
			Stars:       "27,890",
			Topics:      []string{"llm", "rag", "token-optimization", "prompt-engineering", "python"},
			ReadmeNotes: []string{
				"压缩流程先识别输入中的结构化字段、错误栈、指标和实体，再按任务目标保留高价值信息。",
				"项目提供 Python SDK、HTTP Proxy 和 MCP Server 三种接入方式，便于在不同 AI 应用中快速试用。",
				"它解决的是上下文窗口有限和成本高的问题，尤其适合日志分析、客服知识库和代码检索场景。",
			},
		},
		{
			Name:        "openai/openai-cookbook",
			URL:         "https://github.com/openai/openai-cookbook",
			Description: "一个 Agent Skill 注册、发现和编排平台，让团队把可复用能力发布成标准化技能。",
			Language:    "Rust",
			Stars:       "11,760",
			Topics:      []string{"agent", "skill", "workflow", "automation", "rust"},
			ReadmeNotes: []string{
				"平台把技能定义、输入输出 schema、权限和运行入口统一建模，Agent 可以根据任务自动检索和调用技能。",
				"调度层会根据技能声明的能力标签、依赖和执行成本选择合适工具，并记录每次调用结果用于评估。",
				"适合展示 skill、agent、workflow 这类主题在热榜详情页中的呈现效果。",
			},
		},
		{
			Name:        "anthropics/anthropic-cookbook",
			URL:         "https://github.com/anthropics/anthropic-cookbook",
			Description: "面向 RAG 应用的评测看板，跟踪召回质量、答案忠实度、延迟和成本。",
			Language:    "Vue",
			Stars:       "9,604",
			Topics:      []string{"rag", "evaluation", "llm", "dashboard", "vue"},
			ReadmeNotes: []string{
				"项目通过离线测试集和线上反馈同时评估 RAG 链路，覆盖检索、重排、生成和引用校验。",
				"核心指标包括 Recall@K、Faithfulness、Answer Relevance、平均延迟和单次请求成本。",
				"它解决团队上线 RAG 后难以持续观测质量的问题，适合做 AI 应用治理和回归测试。",
			},
		},
	}
}
