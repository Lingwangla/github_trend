package main

import (
	"log"
	"os"
	"time"
)

// main 抓取 GitHub Trending 并同时更新 Markdown 榜单与 GitHub Pages 首页。
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

	if err := writeTrendingMarkdown(repos, updatedAt); err != nil {
		log.Fatalf("write trending markdown failed: %v", err)
	}

	if err := archivePreviousIndex(updatedAt); err != nil {
		log.Fatalf("archive previous index failed: %v", err)
	}

	if err := writeIndexHTML(repos, updatedAt); err != nil {
		log.Fatalf("write index HTML failed: %v", err)
	}
}
