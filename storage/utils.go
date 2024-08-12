package storage

import (
	"fmt"
	"github.com/nikola-susa/pigeon-box/model"
)

func StringSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	}
	if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/1024/1024)
	}
	return fmt.Sprintf("%.2f GB", float64(size)/1024/1024/1024)
}

func isImage(contentType string) bool {
	var imageContentTypes = map[string]bool{
		"image/png":     true,
		"image/jpeg":    true,
		"image/gif":     true,
		"image/webp":    true,
		"image/svg+xml": true,
		"image/avif":    true,
	}

	return imageContentTypes[contentType]
}

func IsPreview(file model.File) bool {
	shouldPreview := false
	if isImage(*file.ContentType) && *file.Size < 1*1024*1024 {
		shouldPreview = true
	}

	return shouldPreview
}
