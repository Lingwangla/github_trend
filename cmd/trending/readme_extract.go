package main

import (
	"strings"
	"unicode"
)

// extractReadmeNotes 从 README 中抽取适合作为项目介绍的关键段落。
func extractReadmeNotes(readme string) []string {
	notes := make([]string, 0, 3)
	for _, line := range strings.Split(readme, "\n") {
		line = cleanMarkdownLine(line)
		if !isUsefulReadmeLine(line) {
			continue
		}

		notes = append(notes, trimLongText(line, 260))
		if len(notes) >= 3 {
			break
		}
	}

	return notes
}

// cleanMarkdownLine 清理 README 单行文本中的常见 Markdown 标记。
func cleanMarkdownLine(line string) string {
	line = strings.TrimSpace(line)
	line = strings.TrimLeft(line, "#>-* ")
	line = strings.TrimSpace(line)
	line = strings.ReplaceAll(line, "`", "")
	line = strings.ReplaceAll(line, "**", "")
	line = strings.ReplaceAll(line, "__", "")
	return line
}

// isUsefulReadmeLine 判断 README 行是否适合作为简介信息。
func isUsefulReadmeLine(line string) bool {
	if len(line) < 50 {
		return false
	}

	lowerLine := strings.ToLower(line)
	skippedPrefixes := []string{
		"!",
		"<",
		"[!",
		"installation",
		"install",
		"usage",
		"quick start",
		"license",
		"contributing",
		"table of contents",
		"import ",
		"const ",
		"let ",
		"var ",
		"func ",
		"type ",
	}
	for _, prefix := range skippedPrefixes {
		if strings.HasPrefix(lowerLine, prefix) {
			return false
		}
	}
	if strings.Contains(lowerLine, "<img") || strings.Contains(lowerLine, "<picture") {
		return false
	}
	if strings.Contains(line, "{") || strings.Contains(line, "}") || strings.Contains(line, ";") {
		return false
	}
	if textDensity(line) < 0.35 {
		return false
	}

	return strings.Contains(line, " ") && !strings.Contains(lowerLine, "http://") && !strings.Contains(lowerLine, "https://")
}

// textDensity 计算文本中有效字母和数字的占比，用于过滤 ASCII art 等噪声。
func textDensity(line string) float64 {
	total := 0
	textChars := 0
	for _, r := range line {
		if unicode.IsSpace(r) {
			continue
		}
		total++
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			textChars++
		}
	}
	if total == 0 {
		return 0
	}

	return float64(textChars) / float64(total)
}

// trimLongText 将长文本截断到指定长度。
func trimLongText(text string, limit int) string {
	if len([]rune(text)) <= limit {
		return text
	}

	runes := []rune(text)
	return string(runes[:limit]) + "..."
}
