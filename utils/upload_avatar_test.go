package utils

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUploadAvatar(t *testing.T) {
	ctx := context.Background()

	image, err := os.Open("../files/file.jpg")
	require.NoError(t, err)

	username := RandomString(6)

	imageURL, err := UploadImage(ctx, image, username)
	require.NoError(t, err)
	require.NotEmpty(t, imageURL)

}
