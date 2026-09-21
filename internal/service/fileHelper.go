package service

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func generateS3Key(userID uint, fileName string) string {

	ext := filepath.Ext(fileName)

	fileID := uuid.New().String()

	now := time.Now()

	return fmt.Sprintf(
		"patients/%d/medical-reports/%d/%02d/%02d/%s%s",
		userID,
		now.Year(),
		now.Month(),
		now.Day(),
		fileID,
		ext,
	)
}