// Package mock provides a Mailer interface to upload files to AWS S3.
package aws

import (
	"log"
	"mime/multipart"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	ResumeBucket string
	Uploader     *s3manager.Uploader
}

func NewS3(resumeBucket string) S3 {
	// AWS credentials are loaded from env automatically
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String("us-east-2")},
	})
	if err != nil {
		log.Fatal("Failed to set up aws session")
	}

	uploader := s3manager.NewUploader(sess)

	return S3{
		ResumeBucket: resumeBucket,
		Uploader:     uploader,
	}
}

func (s3 S3) UploadResume(id int, fh *multipart.FileHeader) (err error) {
	file, err := fh.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s3.Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3.ResumeBucket),
		Key:    aws.String(strconv.Itoa(id) + ".pdf"),
		Body:   file,
	})

	return err
}
