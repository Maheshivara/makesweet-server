package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Save image from context with key = fieldName to server and return the path of the new image
func SaveImageFromContext(ctx *gin.Context, fieldName string) (string, error) {
	image, err := ctx.FormFile(fieldName)
	if err != nil {
		errMsg := fmt.Sprintf("File '%s' not found in form", fieldName)
		return "", errors.New(errMsg)
	}

	allowedMimeTypes := []string{"image/jpeg", "image/png"}
	mimeType, err := getFileType(image)
	if err != nil {
		errMsg := fmt.Sprintf("Fail to assert '%s' extension", fieldName)
		return "", errors.New(errMsg)
	}
	if !slices.Contains(allowedMimeTypes, mimeType) {
		errMsg := fmt.Sprintf("Invalid extension on '%s'", fieldName)
		return "", errors.New(errMsg)
	}

	destPath, err := saveImageFromContext(ctx, image)
	if err != nil {
		errMsg := fmt.Sprintf("Fail to save '%s' in the server", fieldName)
		return "", errors.New(errMsg)
	}
	return destPath, nil
}

func SaveImagesFromContext(ctx *gin.Context, fieldName string) ([]string, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		errMsg := "Fail to load images from form"
		return nil, errors.New(errMsg)
	}

	var imagePaths []string
	images := form.File[fieldName]
	for _, image := range images {
		destPath, err := saveImageFromContext(ctx, image)
		if err != nil {
			errMsg := "Fail to save images in the server"
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
	defer file.Close()

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
		errMsg := "Fail to assert image extension"
		return "", errors.New(errMsg)
	}
	if !slices.Contains(allowedMimeTypes, mimeType) {
		errMsg := "Invalid image extension"
		return "", errors.New(errMsg)
	}

	destFolderPath := os.Getenv("SAVE_IMAGE_FOLDER")
	imageID := uuid.New()
	imageExtension := strings.TrimPrefix(mimeType, "image/")
	imageFileName := fmt.Sprintf("%s.%s", imageID.String(), imageExtension)
	destPath := fmt.Sprintf("%s/%s", destFolderPath, imageFileName)

	err = ctx.SaveUploadedFile(image, destPath)
	if err != nil {
		errMsg := "Fail to save image in the server"
		return "", errors.New(errMsg)
	}

	return destPath, nil
}
