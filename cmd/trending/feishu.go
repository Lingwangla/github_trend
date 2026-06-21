package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// notifyOnlyEnv 是仅发送飞书通知、不重新生成页面的环境变量。
const notifyOnlyEnv = "TREND_NOTIFY_ONLY"

// feishuWebhookURLEnv 是飞书 Webhook 触发器地址环境变量。
const feishuWebhookURLEnv = "FEISHU_WEBHOOK_URL"

// githubRepositoryEnv 是 GitHub Actions 内置仓库名称环境变量。
const githubRepositoryEnv = "GITHUB_REPOSITORY"

// defaultPageURL 是默认 GitHub Pages 页面地址。
const defaultPageURL = "https://lingwangla.github.io/github_trend/"

// notifyFeishuFromReadme 从 Trending Markdown 中提取摘要并发送飞书通知。
func notifyFeishuFromReadme() error {
	webhookURL := strings.TrimSpace(os.Getenv(feishuWebhookURLEnv))
	if webhookURL == "" {
		log.Printf("feishu notification skipped: %s is empty", feishuWebhookURLEnv)
		return nil
	}

	content, err := os.ReadFile(trendingMarkdownPath)
	if err != nil {
		return fmt.Errorf("read trending markdown: %w", err)
	}

	message := buildFeishuMessage(extractReadmeTopRepos(string(content)))
	if err := sendFeishuWebhookText(webhookURL, message); err != nil {
		return fmt.Errorf("send feishu webhook text: %w", err)
	}

	log.Println("feishu notification sent")
	return nil
}

// extractReadmeTopRepos 从 Trending Markdown 中提取 Top 仓库标题。
func extractReadmeTopRepos(content string) []string {
	repos := make([]string, 0, topRepoLimit)
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "## ") {
			continue
		}

		repo := extractMarkdownLinkText(line)
		if repo == "" {
			continue
		}
		repos = append(repos, repo)
		if len(repos) >= topRepoLimit {
			break
		}
	}

	return repos
}

// extractMarkdownLinkText 提取 Markdown 链接中的文本。
func extractMarkdownLinkText(line string) string {
	start := strings.Index(line, "[")
	end := strings.Index(line, "]")
	if start < 0 || end <= start {
		return ""
	}

	return line[start+1 : end]
}

// buildFeishuMessage 构造飞书通知文本。
func buildFeishuMessage(repos []string) string {
	var builder strings.Builder
	builder.WriteString("GitHub AI Daily Trending Top 5 已更新\n")
	builder.WriteString("页面地址：")
	builder.WriteString(currentPageURL())
	builder.WriteString("\n")

	if len(repos) > 0 {
		builder.WriteString("\n今日 Top 5：\n")
		for i, repo := range repos {
			builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, repo))
		}
	}

	return strings.TrimSpace(builder.String())
}

// sendFeishuWebhookText 使用飞书 Webhook 触发器发送文本消息。
func sendFeishuWebhookText(webhookURL string, text string) error {
	client := &http.Client{Timeout: 15 * time.Second}
	body := map[string]any{
		"msg_type": "text",
		"content": map[string]string{
			"text": text,
		},
	}
	_, err := postJSON(client, webhookURL, body)
	return err
}

// currentPageURL 返回当前 GitHub Pages 页面地址。
func currentPageURL() string {
	repository := strings.ToLower(strings.TrimSpace(os.Getenv(githubRepositoryEnv)))
	parts := strings.Split(repository, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return defaultPageURL
	}

	return fmt.Sprintf("https://%s.github.io/%s/", parts[0], parts[1])
}

// postJSON 发送 JSON POST 请求并返回响应体。
func postJSON(client *http.Client, requestURL string, body any) ([]byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected status code: %d body=%s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
