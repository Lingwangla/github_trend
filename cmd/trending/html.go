package main

import (
	"fmt"
	"html"
	"net/url"
	"os"
	"strings"
	"time"
)

// writeIndexHTML 将热榜仓库列表写入 GitHub Pages 首页 HTML。
func writeIndexHTML(repos []repository, updatedAt time.Time) error {
	if err := os.WriteFile(indexPath, []byte(buildIndexHTML(repos, updatedAt)), 0644); err != nil {
		return fmt.Errorf("write index HTML: %w", err)
	}

	return nil
}

// buildIndexHTML 根据热榜仓库列表生成 GitHub Pages 首页 HTML。
func buildIndexHTML(repos []repository, updatedAt time.Time) string {
	var builder strings.Builder

	builder.WriteString(`<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>GitHub `)
	builder.WriteString(html.EscapeString(activeTopicName()))
	builder.WriteString(` Daily Trending Top 5</title>
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
      <h1>GitHub `)
	builder.WriteString(html.EscapeString(activeTopicName()))
	builder.WriteString(` Daily Trending Top 5</h1>
      <p class="updated">更新时间：`)
	builder.WriteString(html.EscapeString(updatedAt.Format(time.RFC3339)))
	builder.WriteString(`</p>
      <p class="updated">`)
	builder.WriteString(html.EscapeString(buildTopicFilterDescription()))
	builder.WriteString(`</p>
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
