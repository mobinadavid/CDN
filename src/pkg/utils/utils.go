package utils

import (
	"cdn/src/config"
	"errors"
	"fmt"
	"github.com/google/uuid"
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

	// todo: each ip allowance in 24 hour: 50 files, 50 mb

	if len(files) > 2 {
		return errors.New("too many files! Maximum 2 files allowed per request")
	}

	for _, file := range files {

		if file.Size > 2*1024*1024 {
			return errors.New("maximum 2MB files allowed")
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
