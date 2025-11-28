package xhsfetcher

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// NoteContent 笔记内容
type NoteContent struct {
	Title        string   // 笔记标题
	Author       string   // 作者名称
	Content      string   // 笔记正文内容
	Images       []string // 图片URL列表
	LikeCount    string   // 点赞数
	CollectCount string   // 收藏数
	CommentCount string   // 评论数
	PublishTime  string   // 发布时间
	Tags         []string // 标签
	URL          string   // 原始URL
}

// Fetcher 小红书内容抓取器
type Fetcher struct {
	client *http.Client
}

// NewFetcher 创建抓取器
func NewFetcher() *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// 允许重定向
				return nil
			},
		},
	}
}

// ParseShareText 解析小红书分享文本
// 格式: 1 【标题 - 作者 | 小红书 - 你的生活兴趣社区】 😆 随机码 😆 https://链接
// 或者直接是URL
func (f *Fetcher) ParseShareText(shareText string) (title, author, url string) {
	shareText = strings.TrimSpace(shareText)

	// 提取URL
	urlPattern := regexp.MustCompile(`https?://[^\s]+`)
	urlMatch := urlPattern.FindString(shareText)
	if urlMatch != "" {
		url = urlMatch
	}

	// 提取【】中的标题信息
	titlePattern := regexp.MustCompile(`【(.+?)】`)
	titleMatch := titlePattern.FindStringSubmatch(shareText)
	if len(titleMatch) > 1 {
		fullTitle := titleMatch[1]
		// 格式: 标题 - 作者 | 小红书 - 你的生活兴趣社区
		parts := strings.Split(fullTitle, " | ")
		if len(parts) > 0 {
			titleAuthor := parts[0]
			// 分离标题和作者
			taParts := strings.Split(titleAuthor, " - ")
			if len(taParts) >= 2 {
				title = strings.TrimSpace(taParts[0])
				author = strings.TrimSpace(taParts[1])
			} else {
				title = strings.TrimSpace(titleAuthor)
			}
		}
	}

	return title, author, url
}

// FetchNote 抓取笔记内容
func (f *Fetcher) FetchNote(url string) (*NoteContent, error) {
	if url == "" {
		return nil, fmt.Errorf("URL不能为空")
	}

	// 构造请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头，模拟浏览器访问
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	// 发送请求
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析HTML
	return f.parseHTML(string(body), url)
}

// parseHTML 解析HTML内容
func (f *Fetcher) parseHTML(html string, originalURL string) (*NoteContent, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}

	note := &NoteContent{
		URL: originalURL,
	}

	// 提取标题 - 尝试多种选择器
	titleSelectors := []string{
		"meta[name='og:title']",
		"meta[property='og:title']",
		"title",
		".title",
		"#detail-title",
		".note-content .title",
	}
	for _, selector := range titleSelectors {
		if selector == "title" {
			note.Title = strings.TrimSpace(doc.Find(selector).Text())
		} else if strings.HasPrefix(selector, "meta") {
			content, exists := doc.Find(selector).Attr("content")
			if exists && content != "" {
				note.Title = strings.TrimSpace(content)
				break
			}
		} else {
			text := strings.TrimSpace(doc.Find(selector).Text())
			if text != "" {
				note.Title = text
				break
			}
		}
	}

	// 提取描述/正文内容
	descSelectors := []string{
		"meta[name='description']",
		"meta[name='og:description']",
		"meta[property='og:description']",
		".desc",
		".content",
		".note-text",
		"#detail-desc",
	}
	for _, selector := range descSelectors {
		if strings.HasPrefix(selector, "meta") {
			content, exists := doc.Find(selector).Attr("content")
			if exists && content != "" {
				note.Content = strings.TrimSpace(content)
				break
			}
		} else {
			text := strings.TrimSpace(doc.Find(selector).Text())
			if text != "" {
				note.Content = text
				break
			}
		}
	}

	// 提取图片
	doc.Find("meta[property='og:image'], meta[name='og:image']").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("content"); exists && src != "" {
			note.Images = append(note.Images, src)
		}
	})

	// 从script标签中提取JSON数据（小红书通常会在页面中嵌入JSON数据）
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		scriptContent := s.Text()
		// 尝试提取 window.__INITIAL_STATE__ 或类似的数据
		if strings.Contains(scriptContent, "__INITIAL_STATE__") ||
			strings.Contains(scriptContent, "noteDetailMap") {
			// 尝试提取笔记内容
			f.extractFromScript(scriptContent, note)
		}
	})

	// 提取标签
	doc.Find(".tag, .hashtag, a[href*='tag']").Each(func(i int, s *goquery.Selection) {
		tag := strings.TrimSpace(s.Text())
		tag = strings.TrimPrefix(tag, "#")
		if tag != "" && len(tag) < 50 {
			note.Tags = append(note.Tags, tag)
		}
	})

	// 清理标题中的平台后缀
	note.Title = f.cleanTitle(note.Title)

	return note, nil
}

// extractFromScript 从script标签中提取数据
func (f *Fetcher) extractFromScript(script string, note *NoteContent) {
	// 尝试提取笔记描述
	descPattern := regexp.MustCompile(`"desc"\s*:\s*"([^"]+)"`)
	if matches := descPattern.FindStringSubmatch(script); len(matches) > 1 {
		content := f.unescapeUnicode(matches[1])
		if content != "" && (note.Content == "" || len(content) > len(note.Content)) {
			note.Content = content
		}
	}

	// 尝试提取标题
	titlePattern := regexp.MustCompile(`"title"\s*:\s*"([^"]+)"`)
	if matches := titlePattern.FindStringSubmatch(script); len(matches) > 1 {
		title := f.unescapeUnicode(matches[1])
		if title != "" && note.Title == "" {
			note.Title = title
		}
	}

	// 尝试提取作者
	authorPattern := regexp.MustCompile(`"nickname"\s*:\s*"([^"]+)"`)
	if matches := authorPattern.FindStringSubmatch(script); len(matches) > 1 {
		author := f.unescapeUnicode(matches[1])
		if author != "" && note.Author == "" {
			note.Author = author
		}
	}

	// 尝试提取点赞数
	likePattern := regexp.MustCompile(`"likedCount"\s*:\s*"?(\d+)"?`)
	if matches := likePattern.FindStringSubmatch(script); len(matches) > 1 {
		note.LikeCount = matches[1]
	}

	// 尝试提取收藏数
	collectPattern := regexp.MustCompile(`"collectedCount"\s*:\s*"?(\d+)"?`)
	if matches := collectPattern.FindStringSubmatch(script); len(matches) > 1 {
		note.CollectCount = matches[1]
	}
}

// unescapeUnicode 解码Unicode转义字符
func (f *Fetcher) unescapeUnicode(s string) string {
	// 处理 \uXXXX 格式
	result := s
	pattern := regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)
	result = pattern.ReplaceAllStringFunc(result, func(match string) string {
		var r rune
		fmt.Sscanf(match, `\u%04x`, &r)
		return string(r)
	})
	// 处理换行符
	result = strings.ReplaceAll(result, `\n`, "\n")
	result = strings.ReplaceAll(result, `\\n`, "\n")
	return result
}

// cleanTitle 清理标题
func (f *Fetcher) cleanTitle(title string) string {
	// 移除 " - 小红书" 等后缀
	suffixes := []string{
		" - 小红书",
		" | 小红书",
		" - 你的生活兴趣社区",
		"小红书 - 你的生活兴趣社区",
	}
	for _, suffix := range suffixes {
		title = strings.TrimSuffix(title, suffix)
	}
	return strings.TrimSpace(title)
}

// FetchFromShareText 从分享文本抓取内容
func (f *Fetcher) FetchFromShareText(shareText string) (*NoteContent, error) {
	title, author, url := f.ParseShareText(shareText)

	if url == "" {
		return nil, fmt.Errorf("未能从分享文本中提取URL")
	}

	note, err := f.FetchNote(url)
	if err != nil {
		// 如果抓取失败，至少返回从分享文本中解析的信息
		return &NoteContent{
			Title:  title,
			Author: author,
			URL:    url,
		}, err
	}

	// 如果抓取的标题为空，使用分享文本中的标题
	if note.Title == "" && title != "" {
		note.Title = title
	}
	if note.Author == "" && author != "" {
		note.Author = author
	}

	return note, nil
}
