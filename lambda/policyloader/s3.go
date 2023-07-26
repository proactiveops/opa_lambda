package policyloader

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	log "github.com/sirupsen/logrus"
)

// S3PolicyLoader loads policies from S3.
type S3PolicyLoader struct {
	bucketName string
	s3Client   s3iface.S3API
}

// NewS3PolicyLoader creates a new S3PolicyLoader.
func NewS3PolicyLoader(bucketName string) (*S3PolicyLoader, error) {
	config := aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}

	sess, err := session.NewSession(&config)
	if err != nil {
		return nil, err
	}

	s3Client := s3.New(sess)
	return &S3PolicyLoader{
		bucketName: bucketName,
		s3Client:   s3Client,
	}, nil
}

// NewS3PolicyLoaderWithClient creates a new S3PolicyLoader with a custom S3 client.
func NewS3PolicyLoaderWithClient(s3Client s3iface.S3API, bucketName string) *S3PolicyLoader {
	return &S3PolicyLoader{
		bucketName: bucketName,
		s3Client:   s3Client,
	}
}

// LoadPolicy loads a policy from S3.
func (loader *S3PolicyLoader) LoadPolicy(ctx context.Context, policyName string) (string, error) {
	objectKey, err := KeyToFilename(policyName)
	if err != nil {
		return "", err
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(loader.bucketName),
		Key:    aws.String(objectKey),
	}

	result, err := loader.s3Client.GetObjectWithContext(ctx, input)
	if err != nil {
		log.Errorf("failed to get policy %s from S3: %v", policyName, err)
		return "", errors.New("failed to get policy from S3")
	}
	defer result.Body.Close()

	content, err := io.ReadAll(result.Body)
	if err != nil {
		log.Errorf("failed to read policy content from %s: %v", policyName, err)
		return "", errors.New("failed to read policy content from S3")
	}

	return string(content), nil
}
