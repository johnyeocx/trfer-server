package cloud

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func ConnectAWS() (s3Cli *s3.Client) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-2"))
	if err != nil {
	  log.Fatalf("failed to load configuration, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	return s3Client
}

func DeleteObject(s3Cli *s3.Client, key string) (error) {
	var bucket = os.Getenv("BUCKET_NAME")

	_, err := s3Cli.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(key),
	})
	return err
}

func GetImageUploadUrl(s3Cli *s3.Client, key string) (string, error) {
	var bucket = os.Getenv("BUCKET_NAME")

	presignClient := s3.NewPresignClient(s3Cli)
	
	presignedUrl, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),

	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(60 * time.Minute)
	})

	if err != nil {
		log.Printf("Couldn't get a presigned request for upload image")
		return "", err
	}

	return presignedUrl.URL, err
}