package cloudinary

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Cloudinary interface {
	UploadImage(file any, params uploader.UploadParams) (string, string, error)
	DeleteImage(publicID string) error
	GetImageURL(publicID string, transformations ...string) string
	DeleteFolder(folderPath string) error
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

func (c *cloudinry) DeleteFolder(folderPath string) error {
	// The public ID prefix for the folder is the folder name followed by a slash.
	prefix := folderPath + "/"

	// Cloudinary's Admin API is used for bulk operations.
	_, err := c.cld.Admin.DeleteAssetsByPrefix(c.ctx, admin.DeleteAssetsByPrefixParams{
		Prefix: []string{prefix},
	})
	if err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	_, err = c.cld.Admin.DeleteFolder(c.ctx, admin.DeleteFolderParams{
		Folder: folderPath,
	})
	if err != nil {
		return fmt.Errorf("failed to delete empty folder: %w", err)
	}

	return nil
}
