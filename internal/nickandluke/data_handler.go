package nickandluke

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const bucket = "nickandluke-guests"
const key = "v1/guests.csv"
const hashedFile = "hashed.txt"

type dataHandler struct {
	sess *session.Session
}

func DataHandler(sess *session.Session) *dataHandler {
	ret := dataHandler{
		sess: sess,
	}
	return &ret
}

func (dh *dataHandler) Download() error {

	file, err := os.Create(guestFile)
	if err != nil {
		return err
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(dh.sess)

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return err
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	return nil
}

func (dh *dataHandler) Upload() error {

	file, err := os.Open(guestFile)
	if err != nil {
		return err
	}

	defer file.Close()

	uploader := s3manager.NewUploader(dh.sess)

	uploadedOutput, err := uploader.Upload(
		&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   file,
		})
	if err != nil {
		return err
	}

	err = saveHash(uploadedOutput)
	if err != nil {
		return err
	}

	fmt.Println("Uploaded", file.Name(), uploadedOutput.UploadID)
	return nil
}

func saveHash(uploadedInput *s3manager.UploadOutput) error {
	file, err := os.Create(hashedFile)
	if err != nil {
		return err
	}

	defer file.Close()
	file.WriteString(strings.Trim(*uploadedInput.ETag, "\""))
	file.WriteString("\n")
	return nil
}
