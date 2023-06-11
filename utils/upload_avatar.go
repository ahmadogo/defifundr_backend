package utils

import (
	"context"

	"mime/multipart"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func UploadAvatar(ctx context.Context, image multipart.File, username string) (string, error) {
	configs, err := LoadConfig("./../")
	if err != nil {
		return "", err
	}

	cld, _ := cloudinary.NewFromURL(configs.CloudinaryURL)

	// Get the preferred name of the file if its not supplied
	fileName := username
	result, err := cld.Upload.Upload(ctx, image, uploader.UploadParams{
		PublicID: fileName,
		Tags:     strings.Split(",", username),
	})
	if err != nil {
		return "", err
	}

	return result.SecureURL, err

}
