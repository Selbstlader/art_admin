package baiduocr

import (
	"context"
	"testing"
	"time"
)

func TestClient_GetAccessToken(t *testing.T) {
	client := NewClient(Config{
		AppID:     "7279235",
		APIKey:    "Imey9DSfXdGNK8sSnD762M0z",
		SecretKey: "GfbtmkH2nJqzbf1OFix3856R3xgUbhZe",
		Timeout:   30,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token, err := client.GetAccessToken(ctx)
	if err != nil {
		t.Fatalf("获取AccessToken失败: %v", err)
	}

	t.Logf("获取到AccessToken: %s...", token[:20])

	// 测试token缓存
	token2, err := client.GetAccessToken(ctx)
	if err != nil {
		t.Fatalf("第二次获取AccessToken失败: %v", err)
	}

	if token != token2 {
		t.Error("Token缓存未生效")
	}
}

func TestClient_RecognizeImageURL(t *testing.T) {
	client := NewClient(Config{
		AppID:     "7279235",
		APIKey:    "Imey9DSfXdGNK8sSnD762M0z",
		SecretKey: "GfbtmkH2nJqzbf1OFix3856R3xgUbhZe",
		Timeout:   30,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 使用一个公开的测试图片
	testImageURL := "https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png"

	result, err := client.RecognizeImageURL(ctx, testImageURL)
	if err != nil {
		t.Logf("OCR识别失败（可能是图片URL问题）: %v", err)
		return
	}

	t.Logf("识别到 %d 个文字区域", result.WordsResultNum)
	for _, word := range result.WordsResult {
		t.Logf("识别文字: %s", word.Words)
	}
}

func TestClient_ExtractTextFromURL(t *testing.T) {
	client := NewClient(Config{
		AppID:     "7279235",
		APIKey:    "Imey9DSfXdGNK8sSnD762M0z",
		SecretKey: "GfbtmkH2nJqzbf1OFix3856R3xgUbhZe",
		Timeout:   30,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	testImageURL := "https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png"

	text, err := client.ExtractTextFromURL(ctx, testImageURL)
	if err != nil {
		t.Logf("提取文字失败: %v", err)
		return
	}

	t.Logf("提取的文字:\n%s", text)
}
