package clients

import (
	"bytes"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/feiyizhou/base/logger"
)

type s3Client struct {
	s3cli *s3.S3
	sess  *session.Session
}

func NewS3Client(region, endpoint, ak, sk string) *s3Client {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		Credentials:      credentials.NewStaticCredentials(ak, sk, ""),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}))
	return &s3Client{
		sess:  sess,
		s3cli: s3.New(sess),
	}
}

func (cli *s3Client) UploadObject(bucket, fullpath string, data []byte) error {
	uploader := s3manager.NewUploader(cli.sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.Concurrency = 4
		u.LeavePartsOnError = true
	})
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fullpath),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}
	logger.Infof("upload object to s3 success, uploadID: %s", result.UploadID)
	return nil
}

func (cli *s3Client) DownloadObject(bucket, srcFullpath, desFullpath string) error {
	file, err := os.Create(desFullpath)
	if err != nil {
		return err
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(cli.sess, func(u *s3manager.Downloader) {
		u.PartSize = 10 * 1024 * 1024
		u.Concurrency = 5
	})
	numBytes, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcFullpath),
	})
	if err != nil {
		return err
	}
	logger.Infof("download object from s3 %s to %s success, file size: %d\n", srcFullpath, desFullpath, numBytes)
	return nil
}

func (cli *s3Client) ListObject(bucket, path string) (*s3.ListObjectsV2Output, error) {
	output, err := cli.s3cli.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(path),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, err
	}
	return output, err
}

func (cli *s3Client) DeleteObject(bucket, fullpath string) error {
	_, err := cli.s3cli.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fullpath),
	})
	return err
}
