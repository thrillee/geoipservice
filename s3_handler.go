package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// downloadFromS3 downloads the MaxMind database from S3 to a local file
func downloadFromS3(ctx context.Context, cfg *Config) error {
	logger := zerolog.Ctx(ctx)

	// Create AWS configuration
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWSAccessKeyID,
			cfg.AWSSecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)

	// Get the object from S3
	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(cfg.S3Bucket),
		Key:    aws.String(cfg.S3Key),
	})
	if err != nil {
		return fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer result.Body.Close()

	// Create the local file
	file, err := os.Create(cfg.LocalDBPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer file.Close()

	// Copy the S3 object content to the local file
	_, err = io.Copy(file, result.Body)
	if err != nil {
		return fmt.Errorf("failed to write local file: %w", err)
	}

	logger.Info().
		Str("bucket", cfg.S3Bucket).
		Str("key", cfg.S3Key).
		Str("local_path", cfg.LocalDBPath).
		Msg("Successfully downloaded MaxMind database from S3")

	return nil
}
