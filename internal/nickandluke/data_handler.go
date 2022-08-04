package nickandluke

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const bucket = "nickandluke-guests"
const key = "v1/guests.csv"
const cache = "staging/guests.csv"

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

	file, err := os.Create(cache)
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
