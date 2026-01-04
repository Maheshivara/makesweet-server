package handlers

import (
	"fmt"
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
	defer os.Remove(templatePath)

	imagePaths, err := utils.SaveImagesFromContext(ctx, "images")
	if err != nil {

		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(imagePaths)

	defer func() {
		for _, imagePath := range imagePaths {
			os.Remove(imagePath)
		}
	}()

	if len(imagePaths) == 0 {
		ctx.JSON(http.StatusBadRequest, "No images found in form")
		return
	}

	destFolderPath := os.Getenv("SAVE_IMAGE_FOLDER")
	outputID := uuid.New()
	outputFileName := fmt.Sprintf("%s.gif", outputID.String())
	outputPath := fmt.Sprintf("%s/%s", destFolderPath, outputFileName)

	customCreateCommand := utils.NewCommandBuilder().CustomTemplate(templatePath, imagePaths, outputPath)
	err = customCreateCommand.Run()
	if err != nil {
		log.Error("Custom gif make fail.", "err", err)
		ctx.JSON(http.StatusInternalServerError, "Fail to create gif")
		return
	}
	defer os.Remove(outputPath)

	_, err = os.Stat(outputPath)
	if err != nil {
		log.Error("Custom output check fail.", "err", err)
		ctx.JSON(http.StatusInternalServerError, "Fail to create gif")
		return
	}
	ctx.File(outputPath)
}
