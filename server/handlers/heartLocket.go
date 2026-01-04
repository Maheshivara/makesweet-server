package handlers

import (
	"fmt"
	"makesweet/utils"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type heartLocketVariant string

const (
	heartLocketVariantStandard  heartLocketVariant = "standard"
	heartLocketVariantChristmas heartLocketVariant = "christmas"
)

// CreateHeartLocket
//
//	@Summary		Create a heart locket gif
//	@Description	Use image-lef and image-right files from form to make a opening heart locket gif
//	@Tags			Gif
//	@Accept			mpfd
//	@Param			image-left	formData	file	true	"A png or jpg image to left half"
//	@Param			image-right	formData	file	true	"A png or jpg image to right half"
//	@Param			variant		formData	string	false	"Variant of heart locket gif"	Enums(standard,christmas) default(standard)
//	@Produce		json image/gif
//	@Success		200	{file}		binary	"Generated Gif"
//	@Failure		400	{string}	string	"Fail to load image from form"
//	@Failure		500	{string}	string	"Fail to generate gif"
//	@Router			/gif/heart-locket [post]
func CreateHeartLocketGif(ctx *gin.Context) {
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
	variant := heartLocketVariant(ctx.DefaultPostForm("variant", "standard"))
	commandBuilder := utils.NewCommandBuilder()
	var command *exec.Cmd
	switch variant {
	case heartLocketVariantChristmas:
		command = commandBuilder.ChristmasHeartLocket(leftImageFilePath, rightImageFilePath, outputPath)
	default:
		command = commandBuilder.HeartLocket(leftImageFilePath, rightImageFilePath, outputPath)
	}
	err = command.Run()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Fail to generate gif")
		return
	}
	defer os.Remove(outputPath)

	ctx.File(outputPath)
}
