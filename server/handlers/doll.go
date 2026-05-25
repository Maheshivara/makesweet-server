package handlers

import (
	"fmt"
	"makesweet/messages"
	"makesweet/utils"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateHeartLocket
//
//	@Summary		Create a nesting doll gif
//	@Description	Use three images to create a nesting doll gif
//	@Tags			Gif
//	@Accept			mpfd
//	@Param			image-left	formData	file	true	"A png or jpg image to the left doll"
//	@Param			image-mid	formData	file	true	"A png or jpg image to the mid doll"
//	@Param			image-right	formData	file	true	"A png or jpg image to the right doll"
//	@Param			mode		formData	string	false	"Gif mode, can be 'fast' or 'slow', default is 'slow'"	Enums(fast,slow)	default(slow)
//	@Produce		json image/gif
//	@Success		200	{file}		binary	"Generated Gif"
//	@Failure		400	{string}	string	"Fail to load images from formData"
//	@Failure		500	{string}	string	"Fail to generate gif"
//	@Router			/gif/nesting-doll [post]
func CreateDollGif(ctx *gin.Context) {
	leftImageFilePath, err := utils.SaveImageFromContext(ctx, "image-left")
	if err != nil {
		expectedErrMsg := fmt.Sprintf(messages.FailToSaveImageInServer, "image-left")
		if err.Error() == expectedErrMsg {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer func() {
		err := os.Remove(leftImageFilePath)
		if err != nil {
			log.Error(messages.FailToRemoveImageFile, "err", err)
		}
	}()

	midImageFilePath, err := utils.SaveImageFromContext(ctx, "image-mid")
	if err != nil {
		expectedErrMsg := fmt.Sprintf(messages.FailToSaveImageInServer, "image-mid")
		if err.Error() == expectedErrMsg {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer func() {
		err := os.Remove(midImageFilePath)
		if err != nil {
			log.Error(messages.FailToRemoveImageFile, "err", err)
		}
	}()

	rightImageFilePath, err := utils.SaveImageFromContext(ctx, "image-right")
	if err != nil {
		expectedErrMsg := fmt.Sprintf(messages.FailToSaveImageInServer, "image-right")
		if err.Error() == expectedErrMsg {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer func() {
		err := os.Remove(rightImageFilePath)
		if err != nil {
			log.Error(messages.FailToRemoveImageFile, "err", err)
		}
	}()

	destDirPath := os.Getenv("SAVE_IMAGE_FOLDER")
	outputID := uuid.New()
	outputFileName := fmt.Sprintf("%s.gif", outputID.String())
	outputPath := fmt.Sprintf("%s/%s", destDirPath, outputFileName)

	fastMode := ctx.DefaultPostForm("mode", "slow") == "fast"
	dollCreateCommand := utils.NewCommandBuilder(fastMode).Doll(leftImageFilePath, midImageFilePath, rightImageFilePath, outputPath)
	err = dollCreateCommand.Run()
	if err != nil {
		log.Error(messages.FailToGenerateGif, "err", err)
		ctx.JSON(http.StatusInternalServerError, messages.FailToGenerateGif)
		return
	}
	defer func() {
		err := os.Remove(outputPath)
		if err != nil {
			log.Error(messages.FailToRemoveGifFile, "err", err)
		}
	}()

	_, err = os.Stat(outputPath)
	if err != nil {
		log.Error(messages.FailToGenerateGif, "err", err)
		ctx.JSON(http.StatusInternalServerError, messages.FailToGenerateGif)
		return
	}
	ctx.File(outputPath)
}
