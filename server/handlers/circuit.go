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
//	@Summary		Create a circuit board gif
//	@Description	Use image from form to make a zooming circuit board gif
//	@Tags			Gif
//	@Accept			mpfd
//	@Param			image	formData	file	true	"A png or jpg image"
//	@Produce		json image/gif
//	@Success		200	{file}		binary	"Generated Gif"
//	@Failure		400	{string}	string	"Fail to load image from form"
//	@Failure		500	{string}	string	"Fail to generate gif"
//	@Router			/gif/circuit [post]
func CreateCircuitGif(ctx *gin.Context) {
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

	billboardCreateCommand := utils.NewCommandBuilder().Circuit(imageFilePath, outputPath)
	err = billboardCreateCommand.Run()
	if err != nil {
		log.Error("Circuit gif make fail.", "err", err)
		ctx.JSON(http.StatusInternalServerError, "Fail to create gif")
		return
	}
	defer os.Remove(outputPath)

	_, err = os.Stat(outputPath)
	if err != nil {
		log.Error("Circuit output check fail.", "err", err)
		ctx.JSON(http.StatusInternalServerError, "Fail to create gif")
		return
	}
	ctx.File(outputPath)
}
