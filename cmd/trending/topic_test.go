package main

import (
	"reflect"
	"testing"
)

// TestParseTopicKeywordsWithChineseSeparators 验证中文分隔符可以正确拆分关键词。
func TestParseTopicKeywordsWithChineseSeparators(t *testing.T) {
	got := parseTopicKeywords("股票，基金、炒股；量化交易")
	want := []string{"股票", "基金", "炒股", "量化交易"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseTopicKeywords() = %v, want %v", got, want)
	}
}

// TestParseTopicKeywordsNormalizesAndDeduplicates 验证关键词归一化和去重逻辑。
func TestParseTopicKeywordsNormalizesAndDeduplicates(t *testing.T) {
	got := parseTopicKeywords("AI, ai; Agent\nagent")
	want := []string{"ai", "agent"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseTopicKeywords() = %v, want %v", got, want)
	}
}
