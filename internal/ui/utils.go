package ui

import (
	"fmt"
	"time"
)

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
