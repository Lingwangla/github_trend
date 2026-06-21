package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

// trendingURL 是 GitHub Daily Trending 页面地址。
const trendingURL = "https://github.com/trending?since=daily"

// readmePath 是 README 输出文件路径。
const readmePath = "README.md"

// indexPath 是 GitHub Pages 入口 HTML 输出文件路径。
const indexPath = "index.html"

// topRepoLimit 是最终输出的 AI 热榜仓库数量。
const topRepoLimit = 5

// mockModeEnv 是开启本地 mock 数据模式的环境变量。
const mockModeEnv = "GITHUB_TREND_USE_MOCK"

// aiKeywords 是用于筛选 AI 相关仓库的关键词。
var aiKeywords = []string{
	"ai",
	"agent",
	"agents",
	"agentic",
	"llm",
	"llms",
	"skill",
	"skills",
	"mcp",
	"model context protocol",
	"chatgpt",
	"openai",
	"claude",
	"gemini",
	"copilot",
	"deepseek",
	"rag",
	"embedding",
	"embeddings",
	"transformer",
	"diffusion",
	"machine learning",
	"ml",
	"deep learning",
	"neural",
	"inference",
	"prompt",
	"prompts",
}

// repository 表示 GitHub 热榜仓库信息。
type repository struct {
	Name        string
	URL         string
	Description string
	Language    string
	Stars       string
	Topics      []string
	ReadmeNotes []string
}

// githubRepositoryResponse 表示 GitHub 仓库 API 的部分响应。
type githubRepositoryResponse struct {
	Description string   `json:"description"`
	Homepage    string   `json:"homepage"`
	Topics      []string `json:"topics"`
}

// githubReadmeResponse 表示 GitHub README API 的部分响应。
type githubReadmeResponse struct {
	DownloadURL string `json:"download_url"`
}

// main 抓取 GitHub Trending 并同时更新 README 与 GitHub Pages 首页。
func main() {
	if os.Getenv(notifyOnlyEnv) == "1" {
		if err := notifyFeishuFromReadme(); err != nil {
			log.Fatalf("notify feishu failed: %v", err)
		}
		return
	}

	repos, err := loadRepositories()
	if err != nil {
		log.Fatalf("load repositories failed: %v", err)
	}

	updatedAt := time.Now().UTC()

	if err := os.WriteFile(readmePath, []byte(buildReadme(repos, updatedAt)), 0644); err != nil {
		log.Fatalf("write README failed: %v", err)
	}

	if err := os.WriteFile(indexPath, []byte(buildIndexHTML(repos, updatedAt)), 0644); err != nil {
		log.Fatalf("write index HTML failed: %v", err)
	}
}

// loadRepositories 根据运行模式加载仓库数据。
func loadRepositories() ([]repository, error) {
	if os.Getenv(mockModeEnv) == "1" {
		return mockRepositories(), nil
	}

	return fetchTrending()
}

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

// fetchTrending 抓取 GitHub Trending 页面并解析仓库列表。
func fetchTrending() ([]repository, error) {
	req, err := http.NewRequest(http.MethodGet, trendingURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", "github-trend-pages-bot")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request trending page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse document: %w", err)
	}

	repos := make([]repository, 0, topRepoLimit)
	doc.Find("article.Box-row").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		name := normalizeRepoName(s.Find("h2 a").Text())
		href, _ := s.Find("h2 a").Attr("href")
		if name == "" || href == "" {
			return true
		}

		repo := repository{
			Name:        name,
			URL:         "https://github.com" + href,
			Description: strings.TrimSpace(s.Find("p").Text()),
			Language:    strings.TrimSpace(s.Find("[itemprop='programmingLanguage']").Text()),
			Stars:       strings.TrimSpace(s.Find("a.Link--muted").First().Text()),
		}
		if !isAIRelatedRepository(repo) {
			return true
		}

		if err := enrichRepository(client, &repo); err != nil {
			log.Printf("enrich %s failed: %v", repo.Name, err)
		}

		repos = append(repos, repo)
		return len(repos) < topRepoLimit
	})

	return repos, nil
}

// enrichRepository 使用 GitHub API 补充仓库主题和 README 摘要信息。
func enrichRepository(client *http.Client, repo *repository) error {
	apiURL := "https://api.github.com/repos/" + repo.Name
	var apiRepo githubRepositoryResponse
	if err := getJSON(client, apiURL, &apiRepo); err != nil {
		return fmt.Errorf("get repository metadata: %w", err)
	}

	if repo.Description == "" {
		repo.Description = apiRepo.Description
	}
	repo.Topics = apiRepo.Topics

	readme, err := fetchReadme(client, repo.Name)
	if err != nil {
		return fmt.Errorf("fetch readme: %w", err)
	}
	repo.ReadmeNotes = extractReadmeNotes(readme)

	return nil
}

// getJSON 读取 JSON 接口并反序列化到目标对象。
func getJSON(client *http.Client, requestURL string, target any) error {
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "github-trend-pages-bot")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

// fetchReadme 获取指定仓库 README 的原始 Markdown 内容。
func fetchReadme(client *http.Client, repoName string) (string, error) {
	var readme githubReadmeResponse
	if err := getJSON(client, "https://api.github.com/repos/"+repoName+"/readme", &readme); err != nil {
		return "", err
	}
	if readme.DownloadURL == "" {
		return "", fmt.Errorf("empty README download URL")
	}

	req, err := http.NewRequest(http.MethodGet, readme.DownloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("create README request: %w", err)
	}
	req.Header.Set("User-Agent", "github-trend-pages-bot")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("download README: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected README status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 256*1024))
	if err != nil {
		return "", fmt.Errorf("read README: %w", err)
	}

	return string(body), nil
}

// extractReadmeNotes 从 README 中抽取适合作为项目介绍的关键段落。
func extractReadmeNotes(readme string) []string {
	notes := make([]string, 0, 3)
	for _, line := range strings.Split(readme, "\n") {
		line = cleanMarkdownLine(line)
		if !isUsefulReadmeLine(line) {
			continue
		}

		notes = append(notes, trimLongText(line, 260))
		if len(notes) >= 3 {
			break
		}
	}

	return notes
}

// cleanMarkdownLine 清理 README 单行文本中的常见 Markdown 标记。
func cleanMarkdownLine(line string) string {
	line = strings.TrimSpace(line)
	line = strings.TrimLeft(line, "#>-* ")
	line = strings.TrimSpace(line)
	line = strings.ReplaceAll(line, "`", "")
	line = strings.ReplaceAll(line, "**", "")
	line = strings.ReplaceAll(line, "__", "")
	return line
}

// isUsefulReadmeLine 判断 README 行是否适合作为简介信息。
func isUsefulReadmeLine(line string) bool {
	if len(line) < 50 {
		return false
	}

	lowerLine := strings.ToLower(line)
	skippedPrefixes := []string{
		"!",
		"<",
		"[!",
		"installation",
		"install",
		"usage",
		"quick start",
		"license",
		"contributing",
		"table of contents",
		"import ",
		"const ",
		"let ",
		"var ",
		"func ",
		"type ",
	}
	for _, prefix := range skippedPrefixes {
		if strings.HasPrefix(lowerLine, prefix) {
			return false
		}
	}
	if strings.Contains(lowerLine, "<img") || strings.Contains(lowerLine, "<picture") {
		return false
	}
	if strings.Contains(line, "{") || strings.Contains(line, "}") || strings.Contains(line, ";") {
		return false
	}
	if textDensity(line) < 0.35 {
		return false
	}

	return strings.Contains(line, " ") && !strings.Contains(lowerLine, "http://") && !strings.Contains(lowerLine, "https://")
}

// textDensity 计算文本中有效字母和数字的占比，用于过滤 ASCII art 等噪声。
func textDensity(line string) float64 {
	total := 0
	textChars := 0
	for _, r := range line {
		if unicode.IsSpace(r) {
			continue
		}
		total++
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			textChars++
		}
	}
	if total == 0 {
		return 0
	}

	return float64(textChars) / float64(total)
}

// trimLongText 将长文本截断到指定长度。
func trimLongText(text string, limit int) string {
	if len([]rune(text)) <= limit {
		return text
	}

	runes := []rune(text)
	return string(runes[:limit]) + "..."
}

// isAIRelatedRepository 判断仓库是否与 AI、Agent、LLM、Skill 等主题相关。
func isAIRelatedRepository(repo repository) bool {
	text := strings.ToLower(strings.Join([]string{
		repo.Name,
		repo.Description,
	}, " "))
	tokens := tokenize(text)

	for _, keyword := range aiKeywords {
		if isShortKeyword(keyword) && tokens[keyword] {
			return true
		}
		if !isShortKeyword(keyword) && strings.Contains(text, keyword) {
			return true
		}
	}

	return false
}

// tokenize 将文本拆分为适合关键词精确匹配的 token 集合。
func tokenize(text string) map[string]bool {
	tokens := make(map[string]bool)
	for _, token := range strings.FieldsFunc(text, func(r rune) bool {
		return (r < 'a' || r > 'z') && (r < '0' || r > '9')
	}) {
		if token != "" {
			tokens[token] = true
		}
	}
	return tokens
}

// isShortKeyword 判断关键词是否需要按 token 精确匹配。
func isShortKeyword(keyword string) bool {
	return len(keyword) <= 3
}

// normalizeRepoName 归一化 GitHub Trending 页面中的仓库名称文本。
func normalizeRepoName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "\n", "")
	return strings.ReplaceAll(name, " ", "")
}

// buildReadme 根据热榜仓库列表生成 README Markdown 内容。
func buildReadme(repos []repository, updatedAt time.Time) string {
	var builder strings.Builder

	builder.WriteString("# GitHub AI Daily Trending Top 5\n\n")
	builder.WriteString(fmt.Sprintf("更新时间：%s\n\n", updatedAt.Format(time.RFC3339)))
	builder.WriteString("筛选范围：仓库名称或描述包含 AI、Agent、LLM、Skill 等相关关键词。\n\n")
	builder.WriteString("网页版本：由 GitHub Pages 自动发布。\n\n")

	for i, repo := range repos {
		builder.WriteString(fmt.Sprintf("## %d. [%s](%s)\n\n", i+1, repo.Name, repo.URL))
		builder.WriteString(fmt.Sprintf("- 语言：%s\n", emptyToDash(repo.Language)))
		builder.WriteString(fmt.Sprintf("- Stars：%s\n", emptyToDash(repo.Stars)))
		builder.WriteString(fmt.Sprintf("- 主题：%s\n", escapeMarkdownCell(formatTopics(repo.Topics))))
		builder.WriteString(fmt.Sprintf("- Star 趋势：\n\n![%s Star History](%s)\n\n", escapeMarkdownCell(repo.Name), buildStarHistoryURL(repo.Name)))
		builder.WriteString(fmt.Sprintf("- 作用 / 解决的问题：%s\n", escapeMarkdownCell(buildPurpose(repo))))
		builder.WriteString("- 适用场景：\n")
		for _, useCase := range buildUseCases(repo) {
			builder.WriteString(fmt.Sprintf("  - %s\n", escapeMarkdownCell(useCase)))
		}
		builder.WriteString("- 架构思想：\n")
		for _, thought := range buildArchitectureThoughts(repo) {
			builder.WriteString(fmt.Sprintf("  - %s\n", escapeMarkdownCell(thought)))
		}
		builder.WriteString("- 原理 / 实现思路：\n")
		for _, principle := range buildPrinciples(repo) {
			builder.WriteString(fmt.Sprintf("  - %s\n", escapeMarkdownCell(principle)))
		}
		builder.WriteString("\n```mermaid\n")
		builder.WriteString(buildArchitectureDiagram(repo))
		builder.WriteString("\n```\n")
		builder.WriteString("\n")
	}

	return builder.String()
}

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
		useCases = append(useCases, "适合围绕 "+formatTopics(repo.Topics)+" 做技术调研、竞品分析或原型验证，因为仓库主题与当前 AI 应用热点高度相关。")
	}

	return useCases
}

// buildArchitectureThoughts 生成仓库架构思想说明。
func buildArchitectureThoughts(repo repository) []string {
	thoughts := []string{
		"它成为热榜的核心原因通常不是单点功能，而是把 LLM、工具、数据和工作流组织成更容易落地的工程结构。",
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
	thoughts = append(thoughts, "它的稀缺性在于把热门 AI 能力包装成可运行、可组合、可观察的工程入口，而不是停留在论文、提示词或孤立 Demo。")

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

// buildIndexHTML 根据热榜仓库列表生成 GitHub Pages 首页 HTML。
func buildIndexHTML(repos []repository, updatedAt time.Time) string {
	var builder strings.Builder

	builder.WriteString(`<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>GitHub AI Daily Trending Top 5</title>
  <style>
    :root {
      color-scheme: light dark;
      --bg: #f6f8fa;
      --card: #ffffff;
      --text: #24292f;
      --muted: #57606a;
      --border: #d0d7de;
      --link: #0969da;
    }
    @media (prefers-color-scheme: dark) {
      :root {
        --bg: #0d1117;
        --card: #161b22;
        --text: #c9d1d9;
        --muted: #8b949e;
        --border: #30363d;
        --link: #58a6ff;
      }
    }
    * {
      box-sizing: border-box;
    }
    body {
      margin: 0;
      background: var(--bg);
      color: var(--text);
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      line-height: 1.5;
    }
    main {
      max-width: 1080px;
      margin: 0 auto;
      padding: 40px 20px;
    }
    header {
      margin-bottom: 24px;
    }
    h1 {
      margin: 0 0 8px;
      font-size: clamp(28px, 5vw, 44px);
      line-height: 1.1;
    }
    .updated {
      color: var(--muted);
      margin: 0;
    }
    .list {
      display: grid;
      gap: 14px;
      padding: 0;
      margin: 0;
      list-style: none;
    }
    .repo {
      background: var(--card);
      border: 1px solid var(--border);
      border-radius: 14px;
      padding: 18px;
    }
    .repo-title {
      display: flex;
      gap: 10px;
      align-items: baseline;
      margin-bottom: 8px;
    }
    .rank {
      color: var(--muted);
      font-variant-numeric: tabular-nums;
    }
    a {
      color: var(--link);
      text-decoration: none;
      font-weight: 650;
    }
    a:hover {
      text-decoration: underline;
    }
    .description {
      margin: 0 0 12px;
      color: var(--text);
    }
    .section {
      margin-top: 14px;
    }
    .section h2 {
      margin: 0 0 6px;
      font-size: 15px;
    }
    .section p {
      margin: 0;
    }
    .section ul {
      margin: 0;
      padding-left: 20px;
    }
    .topics {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      margin-top: 10px;
    }
    .topic {
      border: 1px solid var(--border);
      border-radius: 999px;
      color: var(--muted);
      font-size: 13px;
      padding: 3px 9px;
    }
    .meta {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      color: var(--muted);
      font-size: 14px;
    }
    .star-chart {
      width: 100%;
      margin-top: 12px;
      overflow: hidden;
      border: 1px solid var(--border);
      border-radius: 12px;
      background: #ffffff;
    }
    .star-chart img {
      display: block;
      width: 100%;
      height: auto;
    }
    .mermaid {
      margin-top: 12px;
      padding: 14px;
      overflow-x: auto;
      border: 1px solid var(--border);
      border-radius: 12px;
      background: var(--card);
    }
  </style>
</head>
<body>
  <main>
    <header>
      <h1>GitHub AI Daily Trending Top 5</h1>
      <p class="updated">更新时间：`)
	builder.WriteString(html.EscapeString(updatedAt.Format(time.RFC3339)))
	builder.WriteString(`</p>
      <p class="updated">筛选范围：仓库名称或描述包含 AI、Agent、LLM、Skill 等相关关键词。</p>
    </header>
    <ol class="list">
`)

	for i, repo := range repos {
		builder.WriteString(fmt.Sprintf(`      <li class="repo">
        <div class="repo-title">
          <span class="rank">#%d</span>
          <a href="%s" target="_blank" rel="noopener noreferrer">%s</a>
        </div>
        <div class="meta">
          <span>语言：%s</span>
          <span>Stars：%s</span>
        </div>
        <div class="topics">%s</div>
        <section class="section">
          <h2>Star 趋势</h2>
          <div class="star-chart">
            <img src="%s" alt="%s Star History" loading="lazy">
          </div>
        </section>
        <section class="section">
          <h2>作用 / 解决的问题</h2>
          <p>%s</p>
        </section>
        <section class="section">
          <h2>适用场景</h2>
          %s
        </section>
        <section class="section">
          <h2>架构思想</h2>
          %s
        </section>
        <section class="section">
          <h2>原理 / 实现思路</h2>
          %s
          <div class="mermaid">
%s
          </div>
        </section>
      </li>
`,
			i+1,
			html.EscapeString(repo.URL),
			html.EscapeString(repo.Name),
			html.EscapeString(emptyToDash(repo.Language)),
			html.EscapeString(emptyToDash(repo.Stars)),
			buildTopicsHTML(repo.Topics),
			html.EscapeString(buildStarHistoryURL(repo.Name)),
			html.EscapeString(repo.Name),
			html.EscapeString(buildPurpose(repo)),
			buildListHTML(buildUseCases(repo)),
			buildListHTML(buildArchitectureThoughts(repo)),
			buildListHTML(buildPrinciples(repo)),
			html.EscapeString(buildArchitectureDiagram(repo)),
		))
	}

	builder.WriteString(`    </ol>
  </main>
  <script type="module">
    import mermaid from "https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.esm.min.mjs";
    mermaid.initialize({ startOnLoad: true, theme: "default" });
  </script>
</body>
</html>
`)

	return builder.String()
}

// buildStarHistoryURL 生成仓库 Star History 趋势图地址。
func buildStarHistoryURL(repoName string) string {
	query := url.Values{}
	query.Set("repos", repoName)
	query.Set("type", "Date")
	return "https://api.star-history.com/svg?" + query.Encode()
}

// buildTopicsHTML 生成仓库主题的 HTML 标签。
func buildTopicsHTML(topics []string) string {
	if len(topics) == 0 {
		return `<span class="topic">未公开 topics</span>`
	}

	var builder strings.Builder
	for _, topic := range topics {
		builder.WriteString(`<span class="topic">`)
		builder.WriteString(html.EscapeString(topic))
		builder.WriteString(`</span>`)
	}

	return builder.String()
}

// buildListHTML 生成 HTML 无序列表。
func buildListHTML(items []string) string {
	var builder strings.Builder
	builder.WriteString("<ul>")
	for _, item := range items {
		builder.WriteString("<li>")
		builder.WriteString(html.EscapeString(item))
		builder.WriteString("</li>")
	}
	builder.WriteString("</ul>")
	return builder.String()
}

// escapeMarkdownCell 转义 Markdown 表格单元格中的特殊字符。
func escapeMarkdownCell(value string) string {
	value = strings.ReplaceAll(value, "\n", " ")
	return strings.ReplaceAll(value, "|", "\\|")
}

// emptyToDash 将空字符串转换为占位符。
func emptyToDash(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}
