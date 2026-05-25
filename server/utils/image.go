package utils

import (
	"errors"
	"fmt"
	"makesweet/messages"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Save image from context with key = fieldName to server and return the path of the new image
func SaveImageFromContext(ctx *gin.Context, fieldName string) (string, error) {
	image, err := ctx.FormFile(fieldName)
	if err != nil {
		errMsg := fmt.Sprintf(messages.FileNotFoundInForm, fieldName)
		return "", errors.New(errMsg)
	}

	allowedMimeTypes := []string{"image/jpeg", "image/png"}
	mimeType, err := getFileType(image)
	if err != nil {
		errMsg := fmt.Sprintf(messages.FailToAssertExtension, fieldName)
		return "", errors.New(errMsg)
	}
	if !slices.Contains(allowedMimeTypes, mimeType) {
		errMsg := fmt.Sprintf(messages.InvalidExtensionInForm, fieldName)
		return "", errors.New(errMsg)
	}

	destPath, err := saveImageFromContext(ctx, image)
	if err != nil {
		errMsg := fmt.Sprintf(messages.FailToSaveImageInServer, fieldName)
		return "", errors.New(errMsg)
	}
	return destPath, nil
}

func SaveImagesFromContext(ctx *gin.Context, fieldName string) ([]string, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		errMsg := fmt.Sprintf(messages.FailToLoadImagesFromForm)
		return nil, errors.New(errMsg)
	}

	var imagePaths []string
	images := form.File[fieldName]
	for _, image := range images {
		destPath, err := saveImageFromContext(ctx, image)
		if err != nil {
			errMsg := fmt.Sprintf(messages.FailToSaveImageInServer, fieldName)
			return nil, errors.New(errMsg)
		}
		imagePaths = append(imagePaths, destPath)
	}
	return imagePaths, nil
}

// Get the mimetype from multiform file
func getFileType(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Error(messages.FailToCloseFile, "err", err)
		}
	}()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}

func saveImageFromContext(ctx *gin.Context, image *multipart.FileHeader) (string, error) {
	allowedMimeTypes := []string{"image/jpeg", "image/png"}
	mimeType, err := getFileType(image)
	if err != nil {
		errMsg := fmt.Sprintf(messages.FailToAssertExtension, "image")
		return "", errors.New(errMsg)
	}
	if !slices.Contains(allowedMimeTypes, mimeType) {
		errMsg := fmt.Sprintf(messages.InvalidExtensionInForm, "image")
		return "", errors.New(errMsg)
	}

	destDirPath := os.Getenv("SAVE_IMAGE_FOLDER")
	imageID := uuid.New()
	imageExtension := strings.TrimPrefix(mimeType, "image/")
	imageFileName := fmt.Sprintf("%s.%s", imageID.String(), imageExtension)
	destPath := fmt.Sprintf("%s/%s", destDirPath, imageFileName)

	err = ctx.SaveUploadedFile(image, destPath)
	if err != nil {
		errMsg := fmt.Sprintf(messages.FailToSaveImageInServer, "image")
		return "", errors.New(errMsg)
	}

	return destPath, nil
}
