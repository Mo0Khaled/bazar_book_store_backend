package helpers

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"mime/multipart"
	"os"
)

func credentials() (*cloudinary.Cloudinary, context.Context) {
	cloudName := os.Getenv("IMAGES_CLOUD_NAME")
	apiKey := os.Getenv("IMAGES_API_KEY")
	apiSecret := os.Getenv("IMAGES_API_SECRET")
	cld, _ := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	cld.Config.URL.Secure = true
	ctx := context.Background()
	return cld, ctx
}

func UploadImage(file multipart.File) (string, error) {
	cld, ctx := credentials()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		fmt.Println("error uploading image")
		return "", err
	}

	fmt.Println("****2. Upload an image****\nDelivery URL:", uploadResult.SecureURL)
	return uploadResult.SecureURL, nil
}
