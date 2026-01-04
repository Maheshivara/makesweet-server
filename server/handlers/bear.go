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

// CreateCircuit
//
//	@Summary		Create a flying bear gif
//	@Description	Use image from form to make a flying bear gif
//	@Tags			Gif
//	@Accept			mpfd
//	@Param			image	formData	file	true	"A png or jpg image"
//	@Produce		json image/gif
//	@Success		200	{file}		binary	"Generated Gif"
//	@Failure		400	{string}	string	"Fail to load image from form"
//	@Failure		500	{string}	string	"Fail to generate gif"
//	@Router			/gif/flying-bear [post]
func CreateBearGif(ctx *gin.Context) {
	imageFilePath, err := utils.SaveImageFromContext(ctx, "image")
	if err != nil {
		if err.Error() == "Fail to save 'image' in the server" {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer os.Remove(imageFilePath)

	destFolderPath := os.Getenv("SAVE_IMAGE_FOLDER")
	outputID := uuid.New()
	outputFileName := fmt.Sprintf("%s.gif", outputID.String())
	outputPath := fmt.Sprintf("%s/%s", destFolderPath, outputFileName)

	billboardCreateCommand := utils.NewCommandBuilder().Bear(imageFilePath, outputPath)
	err = billboardCreateCommand.Run()
	if err != nil {
		log.Error("Bear gif make fail.", "err", err)
		ctx.JSON(http.StatusInternalServerError, "Fail to create gif")
		return
	}
	defer os.Remove(outputPath)

	_, err = os.Stat(outputPath)
	if err != nil {
		log.Error("Bear output check fail.", "err", err)
		ctx.JSON(http.StatusInternalServerError, "Fail to create gif")
		return
	}
	ctx.File(outputPath)
}
