package handlers

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func SaveImage(file multipart.File, fileHeader *multipart.FileHeader, dirPath string) (string, error) {
	defer file.Close()

	extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if extension == "" || !isAllowedImageExtension(extension[1:]) {
		return "", errors.New("invalid image extension")
	}

	if fileHeader.Size > 20*1024*1024 {
		return "", errors.New("image too large")
	}

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", err
	}

	fileName := filepath.Base(fileHeader.Filename)
	imagePath := filepath.Join(dirPath, fileName)

	outputFile, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	if _, err := io.Copy(outputFile, file); err != nil {
		return "", err
	}

	return imagePath, nil
}

func isAllowedImageExtension(extension string) bool {
	allowedExtensions := []string{"png", "jpeg", "gif", "webp", "jpg"}
	for _, allowedExt := range allowedExtensions {
		if extension == allowedExt {
			return true
		}
	}
	return false
}
