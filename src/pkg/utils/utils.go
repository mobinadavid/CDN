package utils

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func ValidateFiles(files []*multipart.FileHeader) error {
	log.Println(files)
	allowedExts := []string{".jpg", ".jpeg", ".png", ".pdf", "zip", "rar", "docx", "doc", "csv", "xlsx", "mkv", "mp4"}

	for _, file := range files {
		ext := strings.ToLower(strings.TrimSpace(filepath.Ext(file.Filename)))

		// Check if file extension is allowed
		valid := false
		for _, allowedExt := range allowedExts {
			if ext == allowedExt {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid file format: %s", file.Filename)
		}
	}

	return nil
}

func GenerateUUIDFileName(originalFileName string) string {
	ext := filepath.Ext(originalFileName)
	uuidFileName := uuid.NewString() + ext
	return uuidFileName
}
