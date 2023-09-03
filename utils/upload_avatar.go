package utils

import (
	"context"

	"mime/multipart"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func UploadImage(ctx context.Context, image multipart.File, username string) (string, error) {
	configs, err := LoadConfig("./../")
	if err != nil {
		return "", err
	}

	cld, err := cloudinary.NewFromURL(configs.CloudinaryURL)

	if err != nil {
		return "", err
	}

	// Get the preferred name of the file if its not supplied
	fileName := username
	result, err := cld.Upload.Upload(ctx, image, uploader.UploadParams{
		PublicID: fileName,
		Tags:     strings.Split(",", username),
	})
	if err != nil {
		return "", err
	}

	return result.URL, err

}
