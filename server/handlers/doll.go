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

// CreateHeartLocket
//
//	@Summary		Create a nesting doll gif
//	@Description	Use three images to create a nesting doll gif
//	@Tags			Gif
//	@Accept			mpfd
//	@Param			image-left	formData	file	true	"A png or jpg image to the left doll"
//	@Param			image-mid	formData	file	true	"A png or jpg image to the mid doll"
//	@Param			image-right	formData	file	true	"A png or jpg image to the right doll"
//	@Produce		json image/gif
//	@Success		200	{file}		binary	"Generated Gif"
//	@Failure		400	{string}	string	"Fail to load images from formData"
//	@Failure		500	{string}	string	"Fail to generate gif"
//	@Router			/gif/nesting-doll [post]
func CreateDollGif(ctx *gin.Context) {
	leftImageFilePath, err := utils.SaveImageFromContext(ctx, "image-left")
	if err != nil {
		if err.Error() == "Fail to save 'image-left' in the server" {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer os.Remove(leftImageFilePath)

	midImageFilePath, err := utils.SaveImageFromContext(ctx, "image-mid")
	if err != nil {
		if err.Error() == "Fail to save 'image-mid' in the server" {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer os.Remove(midImageFilePath)

	rightImageFilePath, err := utils.SaveImageFromContext(ctx, "image-right")
	if err != nil {
		if err.Error() == "Fail to save 'image-right' in the server" {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer os.Remove(rightImageFilePath)

	destFolderPath := os.Getenv("SAVE_IMAGE_FOLDER")
	outputID := uuid.New()
	outputFileName := fmt.Sprintf("%s.gif", outputID.String())
	outputPath := fmt.Sprintf("%s/%s", destFolderPath, outputFileName)

	heartLocketCreateCommand := utils.NewCommandBuilder().Doll(leftImageFilePath, midImageFilePath, rightImageFilePath, outputPath)
	err = heartLocketCreateCommand.Run()
	if err != nil {
		log.Error("Nesting doll gif make fail.", "err", err)
		ctx.JSON(http.StatusInternalServerError, "Fail to create gif")
		return
	}
	defer os.Remove(outputPath)

	_, err = os.Stat(outputPath)
	if err != nil {
		log.Error("Nesting doll output check fail.", "err", err)
		ctx.JSON(http.StatusInternalServerError, "Fail to create gif")
		return
	}
	ctx.File(outputPath)
}
