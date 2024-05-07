package utils

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
)

func GenerateUUIDFileName(originalFileName string) string {
	ext := filepath.Ext(originalFileName)
	uuidFileName := uuid.NewString() + ext
	return uuidFileName
}

func ValidateFiles(files []*multipart.FileHeader) error {
	allowedExts := []string{".jpg", ".jpeg", ".png", ".pdf", "zip", "rar", "docx", "doc", "csv", "xlsx", "mkv", "mp4"}

	for _, file := range files {
		log.Println(file)

		ext := strings.ToLower(strings.TrimSpace(filepath.Ext(file.Filename)))

		log.Println(ext)

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
