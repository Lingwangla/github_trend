package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// archivePreviousIndex 将生成前的 index.html 归档为前一天日期命名的历史文件。
func archivePreviousIndex(updatedAt time.Time) error {
	content, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read current index: %w", err)
	}

	if err := os.MkdirAll(historyDailyTrendingDir, 0755); err != nil {
		return fmt.Errorf("create history directory: %w", err)
	}

	historyFileName := updatedAt.AddDate(0, 0, -1).Format("2006-01-02") + ".index"
	historyPath := filepath.Join(historyDailyTrendingDir, historyFileName)
	if err := os.WriteFile(historyPath, content, 0644); err != nil {
		return fmt.Errorf("write history index: %w", err)
	}

	return nil
}
