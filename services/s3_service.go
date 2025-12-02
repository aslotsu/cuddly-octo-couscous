package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Service struct {
	client *s3.Client
	bucket string
}

func NewS3Service() (*S3Service, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1" // default region
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %v", err)
	}

	bucket := os.Getenv("S3_BUCKET_NAME")
	if bucket == "" {
		return nil, fmt.Errorf("S3_BUCKET_NAME environment variable not set")
	}

	client := s3.NewFromConfig(cfg)

	return &S3Service{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *S3Service) UploadImage(file multipart.File, fileSize int64, originalFilename string) (string, string, error) {
	// Validate file type
	ext := strings.ToLower(filepath.Ext(originalFilename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		return "", "", fmt.Errorf("invalid file type: %s", ext)
	}

	// Validate file size (5MB limit)
	if fileSize > 5*1024*1024 {
		return "", "", fmt.Errorf("file too large: maximum 5MB allowed")
	}

	// Generate a unique filename
	fileID := uuid.New().String()
	fileName := fmt.Sprintf("blogs/%s_%s", fileID, originalFilename)

	// Read the file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file: %v", err)
	}

	// Optimize image if needed (for JPEG and PNG)
	fileBuffer := bytes.NewReader(fileBytes)
	optimizedBytes, err := s.optimizeImage(fileBuffer, ext)
	if err != nil {
		// If optimization fails, use the original bytes
		log.Printf("Image optimization failed: %v, using original", err)
		optimizedBytes = fileBytes
	}

	// Upload to S3
	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(optimizedBytes),
		ContentType: aws.String(getContentType(ext)),
		ACL:         "public-read", // Make the image publicly accessible
	})

	if err != nil {
		return "", "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	// Generate public URL
	publicURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, os.Getenv("AWS_REGION"), fileName)

	return fileName, publicURL, nil
}

func (s *S3Service) DeleteImage(key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete from S3: %v", err)
	}

	return nil
}

func (s *S3Service) GetPresignedURL(key string, duration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)
	presignResult, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(duration))

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %v", err)
	}

	return presignResult.URL, nil
}

func getContentType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "image/jpeg" // default to jpeg
	}
}

func (s *S3Service) optimizeImage(reader io.Reader, ext string) ([]byte, error) {
	var img image.Image
	var err error

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(reader)
	case ".png":
		img, err = png.Decode(reader)
	default:
		// For unsupported formats, return original
		return nil, fmt.Errorf("unsupported image format for optimization: %s", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	// Create a buffer to write the optimized image
	var buf bytes.Buffer

	// Encode with quality settings
	switch ext {
	case ".jpg", ".jpeg":
		options := &jpeg.Options{Quality: 80} // 80% quality
		err = jpeg.Encode(&buf, img, options)
	case ".png":
		err = png.Encode(&buf, img) // PNG is lossless but we can optimize it differently if needed
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode optimized image: %v", err)
	}

	return buf.Bytes(), nil
}