package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SaveTemplateFromContext(ctx *gin.Context, fieldName string) (string, error) {
	template, err := ctx.FormFile(fieldName)
	if err != nil {
		errMsg := fmt.Sprintf("File '%s' not found in form", fieldName)
		return "", errors.New(errMsg)
	}

	allowedMimeType := "application/zip"
	mimeType, err := getFileType(template)
	if err != nil {
		errMsg := fmt.Sprintf("Fail to assert '%s' extension", fieldName)
		return "", errors.New(errMsg)
	}
	if mimeType != allowedMimeType {
		errMsg := fmt.Sprintf("Invalid extension on '%s'", fieldName)
		return "", errors.New(errMsg)
	}

	destFolderPath := os.Getenv("SAVE_TEMPLATE_FOLDER")
	templateID := uuid.New()
	templateExtension := "zip"
	templateFileName := fmt.Sprintf("%s.%s", templateID.String(), templateExtension)
	destPath := fmt.Sprintf("%s/%s", destFolderPath, templateFileName)

	err = ctx.SaveUploadedFile(template, destPath)
	if err != nil {
		errMsg := fmt.Sprintf("Fail to save '%s' in the server", fieldName)
		return "", errors.New(errMsg)
	}
	return destPath, nil
}
