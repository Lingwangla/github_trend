package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// writeTrendingMarkdown 将热榜仓库列表写入 Trending Markdown 文件。
func writeTrendingMarkdown(repos []repository, updatedAt time.Time) error {
	if err := os.WriteFile(trendingMarkdownPath, []byte(buildReadme(repos, updatedAt)), 0644); err != nil {
		return fmt.Errorf("write trending markdown: %w", err)
	}

	return nil
}

// buildReadme 根据热榜仓库列表生成 README Markdown 内容。
func buildReadme(repos []repository, updatedAt time.Time) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("# GitHub %s Daily Trending Top 5\n\n", activeTopicName()))
	builder.WriteString(fmt.Sprintf("更新时间：%s\n\n", updatedAt.Format(time.RFC3339)))
	builder.WriteString(buildTopicFilterDescription())
	builder.WriteString("\n\n")
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

// escapeMarkdownCell 转义 Markdown 表格单元格中的特殊字符。
func escapeMarkdownCell(value string) string {
	value = strings.ReplaceAll(value, "\n", " ")
	return strings.ReplaceAll(value, "|", "\\|")
}
