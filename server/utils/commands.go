package utils

import (
	"os/exec"
)

type commandBuilder struct{}
type CommandBuilder interface {
	Billboard(imagePath string, outputPath string) *exec.Cmd
	Flag(imagePath string, outputPath string) *exec.Cmd
	HeartLocket(leftImagePath string, rightImagePath string, outputPath string) *exec.Cmd
	ChristmasHeartLocket(leftImagePath string, rightImagePath string, outputPath string) *exec.Cmd
	Circuit(imagePath string, outputPath string) *exec.Cmd
	Bear(imagePath string, outputPath string) *exec.Cmd
	Doll(imageLeftPath string, imageMidPath string, imageRightPath string, outputPath string) *exec.Cmd
	CustomTemplate(templatePath string, imagePaths []string, outputPath string) *exec.Cmd
}

func NewCommandBuilder() CommandBuilder {
	return &commandBuilder{}
}

func (c *commandBuilder) Billboard(imagePath string, outputPath string) *exec.Cmd {
	images := []string{imagePath}
	cmd := c.reanimateUsingTemplate("/makesweet/templates/billboard-cityscape.zip", images, outputPath)

	return cmd
}

func (c *commandBuilder) Flag(imagePath string, outputPath string) *exec.Cmd {
	images := []string{imagePath}
	cmd := c.reanimateUsingTemplate("/makesweet/templates/flag.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) HeartLocket(leftImagePath string, rightImagePath string, outputPath string) *exec.Cmd {
	images := []string{leftImagePath, rightImagePath}
	cmd := c.reanimateUsingTemplate("/makesweet/templates/heart-locket.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) ChristmasHeartLocket(leftImagePath string, rightImagePath string, outputPath string) *exec.Cmd {
	images := []string{leftImagePath, rightImagePath}
	cmd := c.reanimateUsingTemplate("/makesweet/templates/heart-locket-hat.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) Circuit(imagePath string, outputPath string) *exec.Cmd {
	images := []string{imagePath}
	cmd := c.reanimateUsingTemplate("/makesweet/templates/circuit.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) Bear(imagePath string, outputPath string) *exec.Cmd {
	images := []string{imagePath}
	cmd := c.reanimateUsingTemplate("/makesweet/templates/flying-bear.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) Doll(imageLeftPath string, imageMidPath string, imageRightPath string, outputPath string) *exec.Cmd {
	images := []string{imageLeftPath, imageMidPath, imageRightPath}
	cmd := c.reanimateUsingTemplate("/makesweet/templates/nesting-doll.zip", images, outputPath)
	return cmd
}

func (c *commandBuilder) CustomTemplate(templatePath string, imagePaths []string, outputPath string) *exec.Cmd {
	cmd := c.reanimateUsingTemplate(templatePath, imagePaths, outputPath)
	return cmd
}

func (*commandBuilder) reanimateUsingTemplate(templatePath string, imagePaths []string, outputPath string) *exec.Cmd {
	{
		args := []string{
			"/reanimator",
			"--zip",
			templatePath,
			"--in",
		}
		args = append(args, imagePaths...)
		args = append(args, "--gif", outputPath)

		cmd := exec.Command(args[0], args[1:]...)

		return cmd
	}
}
