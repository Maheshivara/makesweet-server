package handlers

import (
	"fmt"
	"makesweet/messages"
	"makesweet/utils"
	"net/http"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
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
//	@Param			variant		formData	string	false	"Variant of heart locket gif"							Enums(standard,christmas)	default(standard)
//	@Param			mode		formData	string	false	"Gif mode, can be 'fast' or 'slow', default is 'slow'"	Enums(fast,slow)			default(slow)
//	@Produce		json image/gif
//	@Success		200	{file}		binary	"Generated Gif"
//	@Failure		400	{string}	string	"Fail to load image from form"
//	@Failure		500	{string}	string	"Fail to generate gif"
//	@Router			/gif/heart-locket [post]
func CreateHeartLocketGif(ctx *gin.Context) {
	fastMode := ctx.DefaultPostForm("mode", "slow") == "fast"
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
	variant := heartLocketVariant(ctx.DefaultPostForm("variant", "standard"))
	commandBuilder := utils.NewCommandBuilder(fastMode)

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

	defer func() {
		err := os.Remove(outputPath)
		if err != nil {
			log.Error(messages.FailToRemoveGifFile, "err", err)
		}
	}()

	ctx.File(outputPath)
}
