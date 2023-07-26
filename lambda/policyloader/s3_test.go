package policyloader_test

import (
	"context"
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"opa_lambda/policyloader"
)

type mockReadCloser struct {
	mock.Mock
}

func (mrc *mockReadCloser) Read(p []byte) (n int, err error) {
	copy(p, "")
	return 0, errors.New("Bad Read")
}

func (mrc *mockReadCloser) Close() error {
	return nil
}

type mockS3Client struct {
	s3iface.S3API
	mock.Mock
}

func (m *mockS3Client) GetObjectWithContext(ctx aws.Context, input *s3.GetObjectInput, opts ...request.Option) (output *s3.GetObjectOutput, err error) {
	args := m.Called(ctx, input)
	if args.Get(0) != nil {
		output = args.Get(0).(*s3.GetObjectOutput)
	}

	if args.Get(1) != nil {
		err = args.Error(1)
	}

	return
}

func TestLoadItemS3(t *testing.T) {
	s3Client := new(mockS3Client)
	loader := policyloader.NewS3PolicyLoaderWithClient(s3Client, "test-bucket")

	policyName := "test-policy"
	policyContent := "package main\n\ndefault allow = false"

	inputObject := &s3.GetObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String(policyName + ".rego"),
	}

	outputObject := &s3.GetObjectOutput{
		Body: ioutil.NopCloser(strings.NewReader(policyContent)),
	}

	s3Client.On("GetObjectWithContext", mock.Anything, inputObject).Return(outputObject, nil)

	content, err := loader.LoadPolicy(context.Background(), policyName)
	assert.NoError(t, err)
	assert.Equal(t, policyContent, content)

	s3Client.AssertExpectations(t)
}

func TestLoadItemS3_Error(t *testing.T) {
	s3Client := new(mockS3Client)
	loader := policyloader.NewS3PolicyLoaderWithClient(s3Client, "test-bucket")

	policyName := "test-policy"

	s3Client.On("GetObjectWithContext", mock.Anything, &s3.GetObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String(policyName + ".rego"),
	}).Return(nil, errors.New("s3 error"))

	_, err := loader.LoadPolicy(context.Background(), policyName)
	assert.Error(t, err)

	s3Client.AssertExpectations(t)
}

func TestLoadItemS3_Empty(t *testing.T) {
	s3Client := new(mockS3Client)
	loader := policyloader.NewS3PolicyLoaderWithClient(s3Client, "test-bucket")

	policyName := "test-policy"

	inputObject := &s3.GetObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String(policyName + ".rego"),
	}

	outputObject := &s3.GetObjectOutput{
		Body: &mockReadCloser{},
	}

	s3Client.On("GetObjectWithContext", mock.Anything, inputObject).Return(outputObject, nil)

	_, err := loader.LoadPolicy(context.Background(), policyName)
	assert.Error(t, err)

	s3Client.AssertExpectations(t)
}
