package utils

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/stretchr/testify/require"
)

func TestUploadAvatar(t *testing.T) {
	configs, err := LoadConfig("./../")
	fmt.Println(err)
	require.NoError(t, err)

	ctx := context.Background()
	cld, _ := cloudinary.NewFromURL(configs.CloudinaryURL)

	image, err := os.Open("../files/file.jpg")
	require.NoError(t, err)

	username := RandomString(6)
	fileName := username

	cld.Upload.Upload(ctx, image, uploader.UploadParams{
		PublicID: fileName,
		Tags:     strings.Split(",", username),
	})
}
