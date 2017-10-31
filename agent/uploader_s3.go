package agent

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var ErrorUploaderS3Conf = fmt.Errorf("access key, access secret key , region & bucket name is not fully provided")

type UploaderS3 struct {
	isConfigured bool
	bucketName   string
	awsConfig    *aws.Config
}

func (u *UploaderS3) Upload(path string, reader io.Reader) error {
	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("Fail to read content; err : %v", err)
	}

	newBuffer := bytes.NewReader(buffer)
	svc := s3.New(session.New(), u.awsConfig)
	fileType := http.DetectContentType(buffer)
	size := int64(len(buffer))

	params := &s3.PutObjectInput{
		Bucket:        aws.String(u.bucketName),
		Key:           aws.String(path),
		Body:          newBuffer,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	_, err = svc.PutObject(params)
	if err != nil {
		return fmt.Errorf("Fail to upload to S3; err : %v", err)
	}
	return nil

}

func (u *UploaderS3) Configure(conf map[string]string) error {
	if u.isConfigured {
		return nil
	}

	if conf["s3_access_key"] == "" || conf["s3_access_secret_key"] == "" || conf["s3_bucket_name"] == "" || conf["s3_region"] == "" {
		return ErrorUploaderS3Conf
	}

	u.bucketName = conf["s3_bucket_name"]

	token := ""
	creds := credentials.NewStaticCredentials(conf["s3_access_key"], conf["s3_access_secret_key"], token)
	_, err := creds.Get()
	if err != nil {
		return fmt.Errorf("Fail to validate AWS Credential, err : %v", err)
	}

	u.awsConfig = aws.NewConfig().WithRegion(conf["s3_region"]).WithCredentials(creds)
	u.isConfigured = true

	return nil
}

func (u *UploaderS3) Configured() bool {
	return u.isConfigured
}
