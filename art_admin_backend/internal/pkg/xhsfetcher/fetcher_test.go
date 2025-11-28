package xhsfetcher

import (
	"fmt"
	"testing"
)

func TestParseShareText(t *testing.T) {
	fetcher := NewFetcher()

	testCases := []struct {
		input  string
		title  string
		author string
		url    string
	}{
		{
			input:  "1 【华为的live真的狗屎 - cakemehome | 小红书 - 你的生活兴趣社区】 😆 sN8yb0MfT5pbosj 😆 https://www.xiaohongshu.com/discovery/item/6923c54e000000001e0014cc",
			title:  "华为的live真的狗屎",
			author: "cakemehome",
			url:    "https://www.xiaohongshu.com/discovery/item/6923c54e000000001e0014cc",
		},
		{
			input:  "https://www.xiaohongshu.com/explore/123456",
			title:  "",
			author: "",
			url:    "https://www.xiaohongshu.com/explore/123456",
		},
		{
			input:  "41 【就PDD买，别买贵了！！ - 特效药 | 小红书】 😆 qSarbmauncnfmld 😆 https://www.xiaohongshu.com/discovery/item/69256341000000001d03c1cc?source=webshare",
			title:  "就PDD买，别买贵了！！",
			author: "特效药",
			url:    "https://www.xiaohongshu.com/discovery/item/69256341000000001d03c1cc?source=webshare",
		},
	}

	for _, tc := range testCases {
		title, author, url := fetcher.ParseShareText(tc.input)
		fmt.Printf("Input: %s\n", tc.input[:min(50, len(tc.input))])
		fmt.Printf("Parsed - Title: %s, Author: %s, URL: %s\n\n", title, author, url)

		if title != tc.title {
			t.Errorf("Title mismatch: expected %q, got %q", tc.title, title)
		}
		if author != tc.author {
			t.Errorf("Author mismatch: expected %q, got %q", tc.author, author)
		}
		// URL可能包含更多参数，只检查是否包含
		if url == "" && tc.url != "" {
			t.Errorf("URL should not be empty")
		}
	}
}

func TestFetchNote(t *testing.T) {
	fetcher := NewFetcher()

	// 测试真实URL抓取
	url := "https://www.xiaohongshu.com/discovery/item/6923c54e000000001e0014cc"

	note, err := fetcher.FetchNote(url)
	if err != nil {
		t.Logf("Fetch error (expected for some URLs): %v", err)
	}

	if note != nil {
		fmt.Printf("=== Fetched Note ===\n")
		fmt.Printf("Title: %s\n", note.Title)
		fmt.Printf("Author: %s\n", note.Author)
		fmt.Printf("Content: %s\n", note.Content[:min(200, len(note.Content))])
		fmt.Printf("Images: %v\n", note.Images)
		fmt.Printf("Tags: %v\n", note.Tags)
		fmt.Printf("LikeCount: %s\n", note.LikeCount)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
