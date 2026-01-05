package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dakshhhhh16/trace/pkg/aws"
	"github.com/dakshhhhh16/trace/pkg/scan"
	"github.com/gin-gonic/gin"
)

type ImageGeneratePayload struct {
	Arch      string `json:"arch"`
	ImageName string `json:"image_name"`
	Registry  string `json:"registry"`
	OrgID     string `json:"org_id"`   // Optional for now
	ImageID   string `json:"image_id"` // Optional for now
}

// ScanHandler handles scan-related HTTP requests
type ScanHandler struct {
	S3Client   aws.BucketBasics
	BucketName string
}

// NewScanHandler creates a new ScanHandler with dependencies
func NewScanHandler(s3Client aws.BucketBasics, bucketName string) *ScanHandler {
	return &ScanHandler{
		S3Client:   s3Client,
		BucketName: bucketName,
	}
}

// GenerateScanManifestVul generates Manifest, scans for vulnerabilities, and uploads both to S3
func (h *ScanHandler) GenerateScanSbomVul(c *gin.Context) {
	var payload ImageGeneratePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use default org/image IDs if not provided
	orgID := payload.OrgID
	if orgID == "" {
		orgID = "default_org"
	}
	imageID := payload.ImageID
	if imageID == "" {
		imageID = "default_image"
	}

	// Get the source and generate Manifest
	src := scan.GetSource(scan.ImageReference(payload.ImageName))
	defer func() {
		if err := src.Close(); err != nil {
			log.Printf("failed to close source: %v", err)
		}
	}()

	// ==================== Manifest Generation ====================
	manifest := scan.GetManifest(src)

	// Save Manifest to file
	manifestFilePath := "manifest.json"
	if err := scan.SaveManifestToFile(manifest, manifestFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to save Manifest: %v", err),
		})
		return
	}
	defer os.Remove(manifestFilePath)

	// ==================== Vulnerability Scan ====================
	vulnFilePath, err := scan.GetAllVulnAndUpload(manifest)
	if err != nil {
		os.Remove(manifestFilePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to generate vulnerabilities: %v", err),
		})
		return
	}

	// ==================== Upload Both Files to S3 ====================
	ctx := context.Background()

	// Upload Manifest
	manifestS3Key := fmt.Sprintf("trace/%s/%s/manifest.json", orgID, imageID)
	err = h.S3Client.UploadFile(ctx, h.BucketName, manifestS3Key, manifestFilePath)
	if err != nil {
		os.Remove(manifestFilePath)
		os.Remove(vulnFilePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to upload Manifest to S3: %v", err),
		})
		return
	}

	// Upload Vulnerabilities
	vulnS3Key := fmt.Sprintf("trace/%s/%s/vulnerabilities.json", orgID, imageID)
	err = h.S3Client.UploadFile(ctx, h.BucketName, vulnS3Key, vulnFilePath)
	if err != nil {
		os.Remove(manifestFilePath)
		os.Remove(vulnFilePath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to upload vulnerabilities to S3: %v", err),
		})
		return
	}

	// ==================== Clean Up Local Files ====================
	if err := os.Remove(manifestFilePath); err != nil {
		log.Printf("Warning: failed to delete local Manifest file %s: %v", manifestFilePath, err)
	}
	if err := os.Remove(vulnFilePath); err != nil {
		log.Printf("Warning: failed to delete local vulnerability file %s: %v", vulnFilePath, err)
	}

	// ==================== Generate Presigned URLs (Optional) ====================
	manifestPresignedURL, _ := h.S3Client.GetPresignedURL(ctx, h.BucketName, manifestS3Key, time.Hour)
	vulnPresignedURL, _ := h.S3Client.GetPresignedURL(ctx, h.BucketName, vulnS3Key, time.Hour)

	c.JSON(http.StatusOK, gin.H{
		"message":    "Manifest and vulnerabilities generated and uploaded successfully",
		"org_id":     orgID,
		"image_id":   imageID,
		"image_name": payload.ImageName,
		"bucket":     h.BucketName,
		"files": gin.H{
			"manifest": gin.H{
				"s3_key":       manifestS3Key,
				"download_url": manifestPresignedURL,
			},
			"vulnerabilities": gin.H{
				"s3_key":       vulnS3Key,
				"download_url": vulnPresignedURL,
			},
		},
	})
}

// GetImageScans lists all scan files for a specific image
func (h *ScanHandler) GetImageScans(c *gin.Context) {
	orgID := c.Param("org_id")
	imageID := c.Param("image_id")

	if orgID == "" || imageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "org_id and image_id are required",
		})
		return
	}

	ctx := context.Background()
	files, err := h.ListImageScans(ctx, orgID, imageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to list scans: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"org_id":   orgID,
		"image_id": imageID,
		"files":    files,
		"count":    len(files),
	})
}

// DeleteImageScansHandler deletes all scan files for a specific image
func (h *ScanHandler) DeleteImageScansHandler(c *gin.Context) {
	orgID := c.Param("org_id")
	imageID := c.Param("image_id")

	if orgID == "" || imageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "org_id and image_id are required",
		})
		return
	}

	ctx := context.Background()
	err := h.DeleteImageScans(ctx, orgID, imageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to delete scans: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Scans deleted successfully",
		"org_id":   orgID,
		"image_id": imageID,
	})
}

// DownloadScanFile downloads a specific scan file from S3
func (h *ScanHandler) DownloadScanFile(c *gin.Context) {
	orgID := c.Param("org_id")
	imageID := c.Param("image_id")
	filename := c.Param("filename")

	if orgID == "" || imageID == "" || filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "org_id, image_id, and filename are required",
		})
		return
	}

	s3Key := fmt.Sprintf("trace/%s/%s/%s", orgID, imageID, filename)
	localFile := fmt.Sprintf("/tmp/%s", filename)

	ctx := context.Background()
	err := h.S3Client.DownloadFile(ctx, h.BucketName, s3Key, localFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to download file: %v", err),
		})
		return
	}
	defer os.Remove(localFile)

	c.File(localFile)
}
