package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

// GetPresignedUploadURL generates a URL that permits the client to perform an S3 PUT upload directly
func GetPresignedUploadURL(c *gin.Context) {
	filename := c.Query("filename")
	filetype := c.Query("filetype")

	if filename == "" || filetype == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename and filetype parameters are required"})
		return
	}

	bucketName := os.Getenv("AWS_BUCKET_NAME")
	awsRegion := os.Getenv("AWS_REGION")

	// Load AWS configurations from env
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to initialize AWS config: %v", err)})
		return
	}

	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	// Create unique object key inside S3 bucket
	objectKey := fmt.Sprintf("mockups/%d-%s", time.Now().Unix(), filename)

	// Configure S3 PUT pre-signed URL options
	presignParams := &s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         &objectKey,
		ContentType: &filetype,
	}

	// Generate a pre-signed URL valid for 30 minutes
	presignedReq, err := presignClient.PresignPutObject(context.TODO(), presignParams, func(opts *s3.PresignOptions) {
		opts.Expires = 30 * time.Minute
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed generating pre-signed upload URL: %v", err)})
		return
	}

	// Public URL where the uploaded mockup can be downloaded or rendered
	downloadURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, awsRegion, objectKey)

	c.JSON(http.StatusOK, gin.H{
		"uploadUrl":   presignedReq.URL,
		"downloadUrl": downloadURL,
	})
}