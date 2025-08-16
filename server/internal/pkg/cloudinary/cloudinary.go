package cloudinary

import (
	"context"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Cloudinary interface {
	UploadImage(file any, params uploader.UploadParams) (string, string, error)
	DeleteImage(publicID string) error
	GetImageURL(publicID string, transformations ...string) string
}

type cloudinry struct {
	ctx context.Context
	cld *cloudinary.Cloudinary
}

func NewCloudinary() Cloudinary {
	cldnry, _ := cloudinary.New()
	cldnry.Config.URL.Secure = true

	return &cloudinry{
		ctx: context.Background(),
		cld: cldnry,
	}
}

// UploadImage uploads an image and returns its secure URL
func (c *cloudinry) UploadImage(file any, params uploader.UploadParams) (string, string, error) {
	resp, err := c.cld.Upload.Upload(c.ctx, file, params)
	if err != nil {
		return "", "", err
	}
	return resp.SecureURL, resp.PublicID, nil
}

// DeleteImage deletes an image by public ID
func (c *cloudinry) DeleteImage(publicID string) error {
	_, err := c.cld.Upload.Destroy(c.ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

// GetImageURL returns the image URL with optional transformations
// transformations should be strings like "w_300,h_300,c_fill"
func (c *cloudinry) GetImageURL(publicID string, transformations ...string) string {
	if publicID == "" {
		return "" // Return empty string for invalid publicID
	}

	// Create a new Cloudinary image resource
	img, err := c.cld.Image(publicID)
	if err != nil {
		return "" // Return empty string if resource creation fails
	}

	// Apply transformations if provided
	if len(transformations) > 0 {
		// Combine transformations into a single string or apply as separate segments
		img.Transformation = strings.Join(transformations, "/")
	}

	// Generate the secure URL
	url, err := img.String()
	if err != nil {
		return "" // Return empty string if URL generation fails
	}

	return url
}
