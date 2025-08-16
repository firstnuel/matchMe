package user

import (
	"match-me/api/middleware"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *UserHandler) UploadUserPhotos(c *gin.Context) {
	// Get user ID from URL parameter
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"details": "User ID must be a valid UUID",
		})
		return
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to parse multipart form",
			"details": err.Error(),
		})
		return
	}

	// Get files from form
	files := form.File["photos"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No photos provided",
			"details": "At least one photo file is required",
		})
		return
	}

	// Convert multipart files to interface{} slice for usecase
	var fileInterfaces []interface{}
	var openFiles []multipart.File

	for _, fileHeader := range files {
		// Validate file type (basic image validation)
		if !isValidImageFile(fileHeader.Filename) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid file type",
				"details": "Only image files (jpg, jpeg, png, gif, webp) are allowed",
			})
			return
		}

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			// Close any previously opened files
			for _, f := range openFiles {
				f.Close()
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Failed to open uploaded file",
				"details": err.Error(),
			})
			return
		}

		openFiles = append(openFiles, file)
		fileInterfaces = append(fileInterfaces, file)
	}

	// Ensure all files are closed after processing
	defer func() {
		for _, file := range openFiles {
			file.Close()
		}
	}()

	// Upload photos via usecase
	photos, err := h.userUsecase.UploadUserPhotos(c.Request.Context(), userID, fileInterfaces)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to upload photos",
			"details": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Photos uploaded successfully",
		"photos":  photos,
		"count":   len(photos),
	})
}

func (h *UserHandler) UploadCurrentUserPhotos(c *gin.Context) {
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(401, gin.H{"error": "User not found in context"})
		return
	}

	// Set the user ID parameter for the existing UploadUserPhotos handler
	c.Params = append(c.Params, gin.Param{Key: "id", Value: user.ID.String()})
	h.UploadUserPhotos(c)
}

// isValidImageFile checks if the file has a valid image extension
func isValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}
