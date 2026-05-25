package messages

const (
	ExecCommandFailed = "command %s failed."

	FailToRemoveTemplateFile  = "fail to remove template file."
	FailToRemoveImageFile     = "fail to remove image file."
	FailToRemoveGifFile       = "fail to remove gif file."
	FailToGenerateGif         = "fail to generate gif."
	FailToGenerateSpecificGif = "fail to generate %s gif."
	FailToSaveImageInServer   = "fail to save '%s' in the server."
	FailToAssertExtension     = "fail to assert '%s' extension."
	FailToLoadImagesFromForm  = "fail to load images from form."
	FailToCloseFile           = "fail to close file."

	FileNotFoundInForm = "file '%s' not found in form"

	InvalidExtensionInForm = "invalid extension on '%s'"

	NoImagesFoundInForm = "no images found in form"
)
