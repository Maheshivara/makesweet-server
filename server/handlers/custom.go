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

// CreateCustom
//
//	@Summary		Create a custom gif
//	@Description	Use images and a template from form to make a custom gif
//	@Tags			Gif
//	@Accept			mpfd
//	@Param			images		formData	[]file	true	"A png or jpg image array"
//	@Param			template	formData	file	true	"A zip template file"
//	@Param			mode		formData	string	false	"Gif mode, can be 'fast' or 'slow', default is 'slow'"	Enums(fast,slow)	default(slow)
//	@Produce		json image/gif
//	@Success		200	{file}		binary	"Generated Gif"
//	@Failure		400	{string}	string	"Fail to load images from form"
//	@Failure		500	{string}	string	"Fail to generate gif"
//	@Router			/gif/custom [post]
func CreateFromCustom(ctx *gin.Context) {
	templatePath, err := utils.SaveTemplateFromContext(ctx, "template")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer func() {
		err := os.Remove(templatePath)
		if err != nil {
			log.Error(messages.FailToRemoveTemplateFile, "err", err)
		}
	}()

	imagePaths, err := utils.SaveImagesFromContext(ctx, "images")
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(imagePaths)

	defer func() {
		for _, imagePath := range imagePaths {
			err := os.Remove(imagePath)
			if err != nil {
				log.Error(messages.FailToRemoveImageFile, "err", err)
			}
		}
	}()

	if len(imagePaths) == 0 {
		ctx.JSON(http.StatusBadRequest, messages.NoImagesFoundInForm)
		return
	}

	destDirPath := os.Getenv("SAVE_IMAGE_FOLDER")
	outputID := uuid.New()
	outputFileName := fmt.Sprintf("%s.gif", outputID.String())
	outputPath := fmt.Sprintf("%s/%s", destDirPath, outputFileName)

	fastMode := ctx.DefaultPostForm("mode", "slow") == "fast"
	customCreateCommand := utils.NewCommandBuilder(fastMode).CustomTemplate(templatePath, imagePaths, outputPath)
	err = customCreateCommand.Run()
	if err != nil {
		customErr := fmt.Sprintf(messages.ExecCommandFailed, "CustomTemplate")
		log.Error(customErr, "err", err)
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
		customErr := fmt.Sprintf(messages.ExecCommandFailed, "CustomTemplate")
		log.Error(customErr, "err", err)
		ctx.JSON(http.StatusInternalServerError, messages.FailToGenerateGif)
		return
	}
	ctx.File(outputPath)
}
