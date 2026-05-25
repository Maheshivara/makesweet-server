package utils

import (
	"os"
	"os/exec"
)

var (
	SaveImageDir    = os.Getenv("SAVE_IMAGE_FOLDER")
	SaveTemplateDir = os.Getenv("SAVE_TEMPLATE_FOLDER")
	BaseTemplateDir = os.Getenv("BASE_TEMPLATE_FOLDER")
	BaseCommand     = os.Getenv("BASE_COMMAND")
)

type commandBuilder struct {
	fastMode bool
}
type CommandBuilder interface {
	Billboard(imagePath string, outputPath string) *exec.Cmd
	Flag(imagePath string, outputPath string) *exec.Cmd
	SkyFlag(imagePath string, outputPath string) *exec.Cmd
	HeartLocket(leftImagePath string, rightImagePath string, outputPath string) *exec.Cmd
	ChristmasHeartLocket(leftImagePath string, rightImagePath string, outputPath string) *exec.Cmd
	Circuit(imagePath string, outputPath string) *exec.Cmd
	Bear(imagePath string, outputPath string) *exec.Cmd
	Doll(imageLeftPath string, imageMidPath string, imageRightPath string, outputPath string) *exec.Cmd
	FortuneCookie(imagePath string, outputPath string) *exec.Cmd
	Guitar(imagePath string, outputPath string) *exec.Cmd
	Laptop(imagePath string, outputPath string) *exec.Cmd
	CustomTemplate(templatePath string, imagePaths []string, outputPath string) *exec.Cmd
}

func NewCommandBuilder(fastMode bool) CommandBuilder {
	return &commandBuilder{
		fastMode: fastMode,
	}
}

func (c *commandBuilder) Billboard(imagePath string, outputPath string) *exec.Cmd {
	images := []string{imagePath}
	cmd := c.reanimateUsingTemplate(BaseTemplateDir+"/billboard-cityscape.zip", images, outputPath)

	return cmd
}

func (c *commandBuilder) Flag(imagePath string, outputPath string) *exec.Cmd {
	images := []string{imagePath}
	cmd := c.reanimateUsingTemplate(BaseTemplateDir+"/flag.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) SkyFlag(imagePath string, outputPath string) *exec.Cmd {
	images := []string{imagePath}
	cmd := c.reanimateUsingTemplate(BaseTemplateDir+"/blue-sky-flag.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) HeartLocket(leftImagePath string, rightImagePath string, outputPath string) *exec.Cmd {
	images := []string{leftImagePath, rightImagePath}
	cmd := c.reanimateUsingTemplate(BaseTemplateDir+"/heart-locket.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) ChristmasHeartLocket(leftImagePath string, rightImagePath string, outputPath string) *exec.Cmd {
	images := []string{leftImagePath, rightImagePath}
	cmd := c.reanimateUsingTemplate(BaseTemplateDir+"/heart-hat.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) Circuit(imagePath string, outputPath string) *exec.Cmd {
	images := []string{imagePath}
	cmd := c.reanimateUsingTemplate(BaseTemplateDir+"/circuit-board.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) Bear(imagePath string, outputPath string) *exec.Cmd {
	images := []string{imagePath}
	cmd := c.reanimateUsingTemplate(BaseTemplateDir+"/bearplane.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) Doll(imageLeftPath string, imageMidPath string, imageRightPath string, outputPath string) *exec.Cmd {
	images := []string{imageLeftPath, imageMidPath, imageRightPath}
	cmd := c.reanimateUsingTemplate(BaseTemplateDir+"/nesting-doll.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) FortuneCookie(imagePath string, outputPath string) *exec.Cmd {
	images := []string{imagePath}
	cmd := c.reanimateUsingTemplate(BaseTemplateDir+"/fortune-cookie.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) Guitar(imagePath string, outputPath string) *exec.Cmd {
	images := []string{imagePath}
	cmd := c.reanimateUsingTemplate(BaseTemplateDir+"/guitar.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) Laptop(imagePath string, outputPath string) *exec.Cmd {
	images := []string{imagePath}
	cmd := c.reanimateUsingTemplate(BaseTemplateDir+"/laptop.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) CustomTemplate(templatePath string, imagePaths []string, outputPath string) *exec.Cmd {
	cmd := c.reanimateUsingTemplate(templatePath, imagePaths, outputPath)
	return cmd
}

func (c *commandBuilder) reanimateUsingTemplate(templatePath string, imagePaths []string, outputPath string) *exec.Cmd {
	{
		command := BaseCommand
		if len(command) == 0 {
			command = "/bin/makesweet-py"
		}
		args := []string{
			"-S",
			templatePath,
		}
		if c.fastMode {
			args = append(args, "--simple")
		} else {
			args = append(args, "--auto-zoom")
		}
		args = append(args, "--inputs")
		args = append(args, imagePaths...)
		args = append(args, "--gif", outputPath)

		cmd := exec.Command(command, args[1:]...)

		return cmd
	}
}
