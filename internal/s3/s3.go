package s3

import (
	"log"
	"mime/multipart"
	"os"

	"github.com/BoilerMake/bm-app/internal/s3/aws"
	"github.com/BoilerMake/bm-app/internal/s3/mock"
)

// An S3 defines an interface uploading files to S3 like services.
type S3 interface {
	UploadResume(id int, fh *multipart.FileHeader) (err error)
}

// NewS3 creates a new S3 based on the environment mode present. In
// development mode it will return a mock s3 instance and in every other mode
// (really just production) it will return an aws s3 instance.
func NewS3() S3 {
	mode, ok := os.LookupEnv("ENV_MODE")
	if !ok {
		log.Fatalf("environment variable not set: %v. Did you update your .env file?", "ENV_MODE")
	}

	if mode == "development" {
		return mock.NewS3()
	} else {
		resumeBucket, ok := os.LookupEnv("S3_BUCKET_RESUMES")
		if !ok {
			log.Fatalf("environment variable not set: %v. Did you update your .env file?", "S3_BUCKET_RESUMES")
		}

		return aws.NewS3(resumeBucket)
	}
}
