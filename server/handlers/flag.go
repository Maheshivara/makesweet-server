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

// CreateFlag
//
//	@Summary		Create a flag gif
//	@Description	Use image from form to make a waving flag gif
//	@Tags			Gif
//	@Accept			mpfd
//	@Param			image	formData	file	true	"A png or jpg image"
//	@Param			mode	formData	string	false	"Gif mode, can be 'fast' or 'slow', default is 'slow'"	Enums(fast,slow)	default(slow)
//	@Produce		json image/gif
//	@Success		200	{file}		binary	"Generated Gif"
//	@Failure		400	{string}	string	"Fail to load image from formData"
//	@Failure		500	{string}	string	"Fail to generate gif"
//	@Router			/gif/flag [post]
func CreateFlagGif(ctx *gin.Context) {
	imageFilePath, err := utils.SaveImageFromContext(ctx, "image")
	if err != nil {
		expectedError := fmt.Sprintf(messages.FailToSaveImageInServer, "image")
		if err.Error() == expectedError {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer func() {
		err := os.Remove(imageFilePath)
		if err != nil {
			log.Error(messages.FailToRemoveImageFile, "err", err)
		}
	}()

	destDirPath := os.Getenv("SAVE_IMAGE_FOLDER")
	outputID := uuid.New()
	outputFileName := fmt.Sprintf("%s.gif", outputID.String())
	outputPath := fmt.Sprintf("%s/%s", destDirPath, outputFileName)

	fastMode := ctx.DefaultPostForm("mode", "slow") == "fast"
	flagCreateCommand := utils.NewCommandBuilder(fastMode).Flag(imageFilePath, outputPath)
	err = flagCreateCommand.Run()
	if err != nil {
		flagErr := fmt.Sprintf(messages.ExecCommandFailed, "Flag")
		log.Error(flagErr, "err", err)
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
		flagErr := fmt.Sprintf(messages.FailToGenerateSpecificGif, "Flag")
		log.Error(flagErr, "err", err)
		ctx.JSON(http.StatusInternalServerError, messages.FailToGenerateGif)
		return
	}
	ctx.File(outputPath)
}
