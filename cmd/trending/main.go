package main

import (
	"log"
	"os"
	"time"
)

// main 抓取 GitHub Trending 并同时更新 Markdown 榜单与 GitHub Pages 首页。
func main() {
	if os.Getenv(notifyOnlyEnv) == "1" {
		log.Printf("running in notify-only mode with %s=1", notifyOnlyEnv)
		if err := notifyFeishuFromReadme(); err != nil {
			log.Fatalf("notify feishu failed: %v", err)
		}
		return
	}

	log.Printf(
		"starting trending generation: topic=%q keywords=%q raw_%s=%q raw_%s=%q mock=%q",
		activeTopicName(),
		activeTopicKeywords(),
		topicKeywordsEnv,
		os.Getenv(topicKeywordsEnv),
		topicExtraKeywordsEnv,
		os.Getenv(topicExtraKeywordsEnv),
		os.Getenv(mockModeEnv),
	)

	repos, err := loadRepositories()
	if err != nil {
		log.Fatalf("load repositories failed: %v", err)
	}
	log.Printf("loaded repositories: count=%d limit=%d", len(repos), topRepoLimit)
	if len(repos) == 0 {
		log.Printf("warning: no repositories matched topic=%q keywords=%q; generated files will contain only header content", activeTopicName(), activeTopicKeywords())
	}

	updatedAt := time.Now().UTC()

	if err := writeTrendingMarkdown(repos, updatedAt); err != nil {
		log.Fatalf("write trending markdown failed: %v", err)
	}

	if err := archivePreviousIndex(updatedAt); err != nil {
		log.Fatalf("archive previous index failed: %v", err)
	}

	if err := writeIndexHTML(repos, updatedAt); err != nil {
		log.Fatalf("write index HTML failed: %v", err)
	}

	log.Printf("trending generation completed: markdown=%s index=%s repos=%d", trendingMarkdownPath, indexPath, len(repos))
}
