package main

// trendingURL 是 GitHub Daily Trending 页面地址。
const trendingURL = "https://github.com/trending?since=daily"

// trendingMarkdownPath 是 Trending Markdown 榜单输出文件路径。
const trendingMarkdownPath = "TRENDING.md"

// indexPath 是 GitHub Pages 入口 HTML 输出文件路径。
const indexPath = "index.html"

// historyDailyTrendingDir 是每日历史热榜页面归档目录。
const historyDailyTrendingDir = "history_daily_trending"

// topRepoLimit 是最终输出的 AI 热榜仓库数量。
const topRepoLimit = 5

// mockModeEnv 是开启本地 mock 数据模式的环境变量。
const mockModeEnv = "TREND_USE_MOCK"

// topicNameEnv 是热榜主题名称环境变量。
const topicNameEnv = "TREND_TOPIC_NAME"

// topicKeywordsEnv 是覆盖默认主题关键词的环境变量，多个关键词用逗号、分号或换行分隔。
const topicKeywordsEnv = "TREND_KEYWORDS"

// topicExtraKeywordsEnv 是追加默认主题关键词的环境变量，多个关键词用逗号、分号或换行分隔。
const topicExtraKeywordsEnv = "TREND_EXTRA_KEYWORDS"

// defaultTopicName 是默认热榜主题名称。
const defaultTopicName = "AI"

// defaultTopicKeywords 是用于筛选默认主题仓库的关键词。
var defaultTopicKeywords = []string{
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
