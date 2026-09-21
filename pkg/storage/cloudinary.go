package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Storage interface {
	Upload(
		ctx context.Context,
		reader io.Reader,
		publicID string,
	) (*UploadResult, error)

	Delete(
		ctx context.Context,
		publicID string,
	) error
}

type CloudinaryStorage struct {
	Client *cloudinary.Cloudinary
}

type UploadResult struct {
	PublicID     string
	SecureURL    string
	ResourceType string
	Format       string
}

func NewCloudinaryStorage(
	cld *cloudinary.Cloudinary,
) *CloudinaryStorage {

	return &CloudinaryStorage{
		Client: cld,
	}
}

func (s *CloudinaryStorage) Upload(
	ctx context.Context,
	reader io.Reader,
	publicID string,
) (*UploadResult, error) {

	result, err := s.Client.Upload.Upload(
		ctx,
		reader,
		uploader.UploadParams{
			PublicID:       publicID,
			ResourceType:   "auto",
			UniqueFilename: api.Bool(false),
			Overwrite:      api.Bool(false),
		},
	)

	if err != nil {
		return nil, fmt.Errorf(
			"cloudinary upload failed: %w",
			err,
		)
	}

	return &UploadResult{
		PublicID:     result.PublicID,
		SecureURL:    result.SecureURL,
		ResourceType: result.ResourceType,
		Format:       result.Format,
	}, nil
}

func (s *CloudinaryStorage) Delete(
	ctx context.Context,
	publicID string,
) error {

	_, err := s.Client.Upload.Destroy(
		ctx,
		uploader.DestroyParams{
			PublicID:     publicID,
			ResourceType: "image",
		},
	)

	if err != nil {
		return fmt.Errorf(
			"cloudinary delete failed: %w",
			err,
		)
	}

	return nil
}