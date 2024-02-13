package utils

import (
	"fmt"
	"log"
	"mime/multipart"

	"github.com/14jasimmtp/CityVibe-Project-Clean-Architecture/pkg/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func CreateSession() *session.Session {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String(cfg.Aws_region),
			Credentials: credentials.NewStaticCredentials(
				cfg.Aws_access,
				cfg.Aws_secret,
				"",
			),
		},
	))
	return sess
}

func UploadImageToS3(file *multipart.FileHeader, sess *session.Session) (string, error) {
	image, err := file.Open()
	if err != nil {
		return "", err
	}
	// fmt.Println("**", sess)
	defer image.Close()
	uploader := s3manager.NewUploader(sess)
	upload, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("vibecity1/product_images/"),
		Key:    aws.String(file.Filename),
		Body:   image,
		ACL:    aws.String("private"),
	})
	if err != nil {
		fmt.Println("eror")
		return "", err
	}
	return upload.Location, nil
}
