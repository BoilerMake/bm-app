// Package mock provides a S3 interface for uploading files. This mock only
// saves files locally.
package mock

import (
	"io"
	"mime/multipart"
	"os"
	"strconv"
)

type S3 struct{}

func NewS3() S3 {
	return S3{}
}

func (s3 S3) UploadResume(id int, fh *multipart.FileHeader) (err error) {
	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	path := "working/resumes/" + strconv.Itoa(id) + ".pdf"
	if _, err := os.Stat("working/resumes/"); os.IsNotExist(err) {
		os.Mkdir("working/resumes/", 0777)
	}

	dest, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)

	return err
}
