package main

import (
	"fmt"
	"os"
	"strings"
)

// isAIRelatedRepository 判断仓库是否与配置的主题关键词相关。
func isAIRelatedRepository(repo repository) bool {
	text := strings.ToLower(strings.Join([]string{
		repo.Name,
		repo.Description,
	}, " "))
	tokens := tokenize(text)

	for _, keyword := range activeTopicKeywords() {
		if isShortKeyword(keyword) && tokens[keyword] {
			return true
		}
		if !isShortKeyword(keyword) && strings.Contains(text, keyword) {
			return true
		}
	}

	return false
}

// activeTopicName 返回当前热榜主题名称。
func activeTopicName() string {
	value := strings.TrimSpace(os.Getenv(topicNameEnv))
	if value == "" {
		return defaultTopicName
	}

	return value
}

// activeTopicKeywords 返回当前热榜主题关键词列表。
func activeTopicKeywords() []string {
	if keywords := parseTopicKeywords(os.Getenv(topicKeywordsEnv)); len(keywords) > 0 {
		return keywords
	}

	keywords := append([]string{}, defaultTopicKeywords...)
	keywords = append(keywords, parseTopicKeywords(os.Getenv(topicExtraKeywordsEnv))...)
	return normalizeKeywords(keywords)
}

// parseTopicKeywords 解析环境变量中的主题关键词。
func parseTopicKeywords(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	keywords := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r' || r == '\t'
	})
	return normalizeKeywords(keywords)
}

// normalizeKeywords 归一化并去重关键词列表。
func normalizeKeywords(keywords []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(keywords))
	for _, keyword := range keywords {
		keyword = strings.ToLower(strings.TrimSpace(keyword))
		if keyword == "" || seen[keyword] {
			continue
		}

		seen[keyword] = true
		result = append(result, keyword)
	}

	return result
}

// buildTopicFilterDescription 生成当前主题筛选说明。
func buildTopicFilterDescription() string {
	return fmt.Sprintf("筛选范围：仓库名称或描述包含 %s 相关关键词。关键词：%s。", activeTopicName(), strings.Join(activeTopicKeywords(), ", "))
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
