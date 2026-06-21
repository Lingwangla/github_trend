# GitHub Trend Pages

GitHub Trend Pages 是一个用 Go 编写的 GitHub Daily Trending 静态榜单生成器。它会抓取 GitHub 每日 Trending 页面，根据可配置的主题关键词筛选 Top 5 仓库，补充仓库 topics 和 README 摘要，并生成 GitHub Pages 可直接发布的静态页面。

当前默认主题是 AI，默认关键词覆盖 AI、Agent、LLM、MCP、RAG、Embedding、Transformer、Prompt 等方向。主题名称和关键词都可以通过环境变量定制。

## 功能

- 抓取 GitHub Daily Trending。
- 按主题关键词筛选目标仓库。
- 通过 GitHub API 补充仓库 topics 和 README 摘要。
- 生成 `TRENDING.md` Markdown 榜单。
- 生成 `index.html` GitHub Pages 首页。
- 更新首页前将旧 `index.html` 归档到 `history_daily_trending/`。
- 支持 GitHub Actions 定时更新、提交产物并部署 GitHub Pages。
- 支持部署完成后发送飞书文本通知。
- 支持 mock 数据模式，便于本地预览页面效果。

## 技术栈

- Go 1.22
- goquery
- GitHub Actions
- GitHub Pages
- 飞书开放平台消息 API

## 目录结构

```text
.
├── .github/workflows/update.yml   # 定时更新、提交生成产物、部署 Pages、发送飞书通知
├── cmd/trending/main.go           # 抓取、筛选、页面生成、历史归档主流程
├── cmd/trending/feishu.go         # 飞书通知逻辑
├── history_daily_trending/        # 历史首页归档目录，文件名格式为 YYYY-MM-DD.index
├── go.mod                         # Go module 和依赖声明
├── go.sum                         # Go 依赖锁定文件
├── README.md                      # 当前项目说明文档
├── TRENDING.md                    # 自动生成的 Trending Markdown 榜单
└── index.html                     # 自动生成的 GitHub Pages 首页
```

## 本地运行

环境要求：

- Go 1.22 或更高版本
- 可以访问 `github.com` 和 `api.github.com`

安装依赖：

```bash
go mod download
```

抓取真实 GitHub Trending 并生成产物：

```bash
go run ./cmd/trending
```

运行后会更新：

- `TRENDING.md`
- `index.html`
- `history_daily_trending/YYYY-MM-DD.index`，仅当运行前已有 `index.html` 时生成

使用 mock 数据本地预览：

```bash
GITHUB_TREND_USE_MOCK=1 go run ./cmd/trending
```

本地预览 GitHub Pages 首页：

```bash
python3 -m http.server 8000
```

然后访问 `http://localhost:8000/index.html`。

## 部署

仓库内置 GitHub Actions 工作流：`.github/workflows/update.yml`。

触发方式：

- 每天 UTC 00:00 自动执行。
- 在 GitHub Actions 页面手动触发 `workflow_dispatch`。

工作流流程：

1. 拉取仓库代码。
2. 安装 Go 1.22。
3. 注入主题相关变量并执行 `go run ./cmd/trending`。
4. 提交生成产物：`TRENDING.md`、`index.html`、`history_daily_trending/`。
5. 上传仓库根目录作为 GitHub Pages artifact。
6. 部署到 GitHub Pages。
7. 部署成功后以通知模式再次执行程序，读取 `TRENDING.md` 并发送飞书通知。

### GitHub Pages

在仓库 `Settings -> Pages` 中，将部署来源设置为 `GitHub Actions`。

工作流已声明 Pages 部署所需权限：

```yaml
permissions:
  contents: write
  pages: write
  id-token: write
```

### GitHub Actions 配置位置

非敏感配置建议放在 `Settings -> Secrets and variables -> Actions -> Variables`。

敏感配置建议放在 `Settings -> Secrets and variables -> Actions -> Secrets`。

## 环境变量

### 榜单生成

| 环境变量 | 是否必填 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `GITHUB_TREND_USE_MOCK` | 否 | 空 | 值为 `1` 时使用内置 mock 数据，不请求 GitHub Trending。 |
| `GITHUB_TREND_TOPIC_NAME` | 否 | `AI` | 榜单主题名称，会展示在 `TRENDING.md` 和 `index.html` 的筛选说明中。 |
| `GITHUB_TREND_KEYWORDS` | 否 | 空 | 覆盖默认主题关键词。配置后只使用该变量中的关键词。 |
| `GITHUB_TREND_EXTRA_KEYWORDS` | 否 | 空 | 在默认主题关键词基础上追加关键词。 |

关键词变量支持用逗号、分号、换行、回车或 tab 分隔，程序会统一转为小写并去重。

覆盖默认关键词示例：

```bash
GITHUB_TREND_TOPIC_NAME=Database \
GITHUB_TREND_KEYWORDS="database,postgres,mysql,redis,sqlite" \
go run ./cmd/trending
```

追加默认关键词示例：

```bash
GITHUB_TREND_EXTRA_KEYWORDS="workflow,automation" go run ./cmd/trending
```

### 飞书通知

| 环境变量 | 是否必填 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `GITHUB_TREND_NOTIFY_ONLY` | 否 | 空 | 值为 `1` 时进入通知模式，只读取 `TRENDING.md` 并发送飞书通知，不重新抓取或生成页面。 |
| `FEISHU_NOTIFY_ENABLED` | 否 | 空 | 值为 `true`、`1` 或 `yes` 时启用飞书通知。 |
| `FEISHU_APP_ID` | 启用通知时必填 | 空 | 飞书自建应用 App ID。 |
| `FEISHU_APP_SECRET` | 启用通知时必填 | 空 | 飞书自建应用 App Secret。 |
| `FEISHU_RECEIVE_ID` | 启用通知时必填 | 空 | 消息接收者 ID。 |
| `FEISHU_RECEIVE_ID_TYPE` | 否 | `email` | 接收者 ID 类型，例如 `email`、`open_id`、`user_id`、`chat_id`。 |
| `FEISHU_PAGE_URL` | 否 | `https://lingwangla.github.io/github_trend/` | 飞书消息中的页面地址。GitHub Actions 中使用 Pages 部署步骤输出的 URL。 |

仅发送飞书通知示例：

```bash
GITHUB_TREND_NOTIFY_ONLY=1 \
FEISHU_NOTIFY_ENABLED=true \
FEISHU_APP_ID=cli_xxx \
FEISHU_APP_SECRET=xxx \
FEISHU_RECEIVE_ID=user@example.com \
go run ./cmd/trending
```

## Actions 配置示例

Variables：

```text
GITHUB_TREND_TOPIC_NAME=AI
GITHUB_TREND_EXTRA_KEYWORDS=workflow,automation
FEISHU_NOTIFY_ENABLED=true
FEISHU_RECEIVE_ID_TYPE=email
```

Secrets：

```text
FEISHU_APP_ID=cli_xxx
FEISHU_APP_SECRET=xxx
FEISHU_RECEIVE_ID=user@example.com
```

## 产物说明

- `README.md` 是项目说明文档，不会被生成器覆盖。
- `TRENDING.md` 是每日 Trending Markdown 榜单。
- `index.html` 是当前 GitHub Pages 首页。
- `history_daily_trending/YYYY-MM-DD.index` 是运行时归档的上一版首页。

## 注意事项

- 当前 GitHub API 请求未配置 token，频繁运行可能触发 GitHub 未认证请求限流。
- GitHub Trending 页面结构变化时，`cmd/trending/main.go` 中的解析选择器可能需要同步调整。
- 历史归档使用当前更新时间的前一天作为文件名，首次运行或不存在旧 `index.html` 时不会生成归档文件。
- 飞书通知依赖自建应用权限、接收者 ID 类型和接收者 ID 配置，失败原因会输出到 Actions 日志。

## 验证

```bash
go test ./...
```
