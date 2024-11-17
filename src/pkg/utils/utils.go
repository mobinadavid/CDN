package utils

import (
	"cdn/src/config"
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
)

func GenerateUUIDFileName(originalFileName string) string {
	ext := filepath.Ext(originalFileName)
	uuidFileName := uuid.NewString() + ext
	return uuidFileName
}

func ValidateFiles(files []*multipart.FileHeader) error {
	fileNumber, err := strconv.Atoi(config.GetInstance().Get("MINIO_MAX_FILES"))
	if err != nil {
		fileNumber = 2
	}
	if len(files) > fileNumber {
		return fmt.Errorf("too many files! Maximum %d files allowed per request", fileNumber)
	}

	for _, file := range files {
		fileSize, err := strconv.Atoi(config.GetInstance().Get("MINIO_FILES_SIZE_MB"))
		if err != nil {
			fileSize = 2
		}
		if file.Size > int64(fileSize)*1024*1024 {
			return fmt.Errorf("maximum file size of %d MB is allowed", fileSize)
		}

		ext := strings.ToLower(strings.TrimSpace(filepath.Ext(file.Filename)))

		// Check if file extension is allowed
		var configs = config.GetInstance()
		allowedExtsEnv := configs.Get("ALLOWED_EXTENSIONS")
		allowedExts := strings.Split(allowedExtsEnv, ",")
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
