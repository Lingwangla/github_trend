package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// notifyOnlyEnv 是仅发送飞书通知、不重新生成页面的环境变量。
const notifyOnlyEnv = "GITHUB_TREND_NOTIFY_ONLY"

// feishuNotifyEnabledEnv 是开启飞书通知的环境变量。
const feishuNotifyEnabledEnv = "FEISHU_NOTIFY_ENABLED"

// feishuAppIDEnv 是飞书应用 App ID 环境变量。
const feishuAppIDEnv = "FEISHU_APP_ID"

// feishuAppSecretEnv 是飞书应用 App Secret 环境变量。
const feishuAppSecretEnv = "FEISHU_APP_SECRET"

// feishuReceiveIDEnv 是飞书消息接收者 ID 环境变量。
const feishuReceiveIDEnv = "FEISHU_RECEIVE_ID"

// feishuReceiveIDTypeEnv 是飞书消息接收者 ID 类型环境变量。
const feishuReceiveIDTypeEnv = "FEISHU_RECEIVE_ID_TYPE"

// feishuPageURLEnv 是飞书通知中的 GitHub Pages 地址环境变量。
const feishuPageURLEnv = "FEISHU_PAGE_URL"

// defaultFeishuReceiveIDType 是默认飞书消息接收者 ID 类型。
const defaultFeishuReceiveIDType = "email"

// defaultPageURL 是默认 GitHub Pages 页面地址。
const defaultPageURL = "https://lingwangla.github.io/github_trend/"

// feishuTenantTokenResponse 表示飞书 tenant_access_token 接口响应。
type feishuTenantTokenResponse struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	TenantAccessToken string `json:"tenant_access_token"`
}

// feishuMessageResponse 表示飞书发送消息接口响应。
type feishuMessageResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// notifyFeishuFromReadme 从 README 中提取摘要并发送飞书通知。
func notifyFeishuFromReadme() error {
	if !isFeishuNotificationEnabled() {
		log.Println("feishu notification skipped: FEISHU_NOTIFY_ENABLED is not true")
		return nil
	}

	content, err := os.ReadFile(readmePath)
	if err != nil {
		return fmt.Errorf("read README: %w", err)
	}

	message := buildFeishuMessage(extractReadmeTopRepos(string(content)))
	if err := sendFeishuText(message); err != nil {
		return fmt.Errorf("send feishu text: %w", err)
	}

	log.Println("feishu notification sent")
	return nil
}

// isFeishuNotificationEnabled 判断是否开启飞书通知。
func isFeishuNotificationEnabled() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(feishuNotifyEnabledEnv)))
	return value == "true" || value == "1" || value == "yes"
}

// extractReadmeTopRepos 从 README 中提取 Top 仓库标题。
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
	builder.WriteString(getEnvDefault(feishuPageURLEnv, defaultPageURL))
	builder.WriteString("\n")

	if len(repos) > 0 {
		builder.WriteString("\n今日 Top 5：\n")
		for i, repo := range repos {
			builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, repo))
		}
	}

	return strings.TrimSpace(builder.String())
}

// sendFeishuText 使用飞书开放平台发送文本消息。
func sendFeishuText(text string) error {
	appID := strings.TrimSpace(os.Getenv(feishuAppIDEnv))
	appSecret := strings.TrimSpace(os.Getenv(feishuAppSecretEnv))
	receiveID := strings.TrimSpace(os.Getenv(feishuReceiveIDEnv))
	receiveIDType := getEnvDefault(feishuReceiveIDTypeEnv, defaultFeishuReceiveIDType)
	if appID == "" || appSecret == "" || receiveID == "" {
		return fmt.Errorf("missing required env: %s, %s or %s", feishuAppIDEnv, feishuAppSecretEnv, feishuReceiveIDEnv)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	token, err := fetchFeishuTenantAccessToken(client, appID, appSecret)
	if err != nil {
		return fmt.Errorf("fetch tenant access token: %w", err)
	}

	return sendFeishuMessage(client, token, receiveIDType, receiveID, text)
}

// fetchFeishuTenantAccessToken 获取飞书 tenant_access_token。
func fetchFeishuTenantAccessToken(client *http.Client, appID string, appSecret string) (string, error) {
	body := map[string]string{
		"app_id":     appID,
		"app_secret": appSecret,
	}

	respBody, err := postJSON(client, "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal", "", body)
	if err != nil {
		return "", err
	}

	var tokenResp feishuTenantTokenResponse
	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}
	if tokenResp.Code != 0 {
		return "", fmt.Errorf("feishu token response code=%d msg=%s", tokenResp.Code, tokenResp.Msg)
	}
	if tokenResp.TenantAccessToken == "" {
		return "", fmt.Errorf("empty tenant access token")
	}

	return tokenResp.TenantAccessToken, nil
}

// sendFeishuMessage 发送飞书文本消息。
func sendFeishuMessage(client *http.Client, token string, receiveIDType string, receiveID string, text string) error {
	content, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return fmt.Errorf("marshal message content: %w", err)
	}

	requestURL := "https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=" + url.QueryEscape(receiveIDType)
	body := map[string]string{
		"receive_id": receiveID,
		"msg_type":   "text",
		"content":    string(content),
	}

	respBody, err := postJSON(client, requestURL, token, body)
	if err != nil {
		return err
	}

	var messageResp feishuMessageResponse
	if err := json.Unmarshal(respBody, &messageResp); err != nil {
		return fmt.Errorf("decode message response: %w", err)
	}
	if messageResp.Code != 0 {
		return fmt.Errorf("feishu message response code=%d msg=%s", messageResp.Code, messageResp.Msg)
	}

	return nil
}

// postJSON 发送 JSON POST 请求并返回响应体。
func postJSON(client *http.Client, requestURL string, bearerToken string, body any) ([]byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

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

// getEnvDefault 获取环境变量，未配置时返回默认值。
func getEnvDefault(key string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}

	return value
}
