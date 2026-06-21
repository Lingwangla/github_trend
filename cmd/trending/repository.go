package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

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

// loadRepositories 根据运行模式加载仓库数据。
func loadRepositories() ([]repository, error) {
	if os.Getenv(mockModeEnv) == "1" {
		return mockRepositories(), nil
	}

	return fetchTrending()
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
