package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func s3exists(key string) (bool, error) {
	_, err := s3service.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s3bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		return true, nil
	}
	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
		return false, nil
	}
	return false, err
}

func s3upload(key string, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	uploader := s3manager.NewUploader(s3session)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	return err
}

func s3download(key string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	downloader := s3manager.NewDownloader(s3session)

	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(s3bucket),
		Key:    aws.String(key),
	})

	return err
}
