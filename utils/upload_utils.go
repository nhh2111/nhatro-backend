package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func ProcessImageUpload(c *gin.Context, file *multipart.FileHeader, folder string, oldImageURL string) (string, error) {
	DeleteOldFile(folder, oldImageURL)

	uploadDir := filepath.Join(".", "uploads", folder)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", errors.New("không thể tạo thư mục lưu trữ")
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return "", errors.New("lỗi khi lưu file ảnh")
	}

	return fmt.Sprintf("/uploads/%s/%s", folder, filename), nil
}

func DeleteOldFile(folder string, oldImageURL string) {
	if oldImageURL == "" {
		return
	}

	fileName := filepath.Base(oldImageURL)
	oldFilePath := filepath.Join(".", "uploads", folder, fileName)

	_ = os.Remove(oldFilePath)
}

func ProcessMultipleImagesUpload(c *gin.Context, files []*multipart.FileHeader, folder string) ([]string, error) {
	var uploadedURLs []string

	uploadDir := filepath.Join(".", "uploads", folder)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, errors.New("không thể tạo thư mục lưu trữ")
	}

	for _, file := range files {
		filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		savePath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(file, savePath); err == nil {
			fileURL := fmt.Sprintf("/uploads/%s/%s", folder, filename)
			uploadedURLs = append(uploadedURLs, fileURL)
		}
	}

	return uploadedURLs, nil
}
