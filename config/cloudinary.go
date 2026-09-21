package config

import (
	//"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

func InitCloudinary() (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.New()
	if err != nil {
		return nil, err
	}

	cld.Config.URL.Secure = true

	return cld, nil
}