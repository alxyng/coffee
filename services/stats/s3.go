package stats

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

type S3StatsOptions struct {
	Bucket     string
	Key        string
	Downloader Downloader
	Uploader   Uploader
}

type S3StatsService struct {
	bucket     *string
	key        *string
	downloader Downloader
	uploader   Uploader
}

type Downloader interface {
	Download() ([]byte, error)
}

type S3Downloader struct {
	bucket   *string
	key      *string
	s3Client *s3.S3
}

func NewS3Downloader(bucket string, key string, s3Client *s3.S3) *S3Downloader {
	return &S3Downloader{
		bucket:   aws.String(bucket),
		key:      aws.String(key),
		s3Client: s3Client,
	}
}

func (d S3Downloader) Download() ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: d.bucket,
		Key:    d.key,
	}

	obj, err := d.s3Client.GetObject(input)
	if err != nil {
		return nil, errors.Wrap(err, "error getting object")
	}
	defer obj.Body.Close()

	data, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading body")
	}

	return data, nil
}

type Uploader interface {
	Upload([]byte) error
}

type S3Uploader struct {
	bucket   *string
	key      *string
	s3Client *s3.S3
}

func NewS3Uploader(bucket string, key string, s3Client *s3.S3) *S3Uploader {
	return &S3Uploader{
		bucket:   aws.String(bucket),
		key:      aws.String(key),
		s3Client: s3Client,
	}
}

func (u S3Uploader) Upload(data []byte) error {
	input := &s3.PutObjectInput{
		Bucket: u.bucket,
		Key:    u.key,
		Body:   aws.ReadSeekCloser(bytes.NewReader(data)),
	}

	_, err := u.s3Client.PutObject(input)
	if err != nil {
		return errors.Wrap(err, "error putting object")
	}

	return nil
}

func NewS3StatsService(options S3StatsOptions) *S3StatsService {
	return &S3StatsService{
		bucket:     aws.String(options.Bucket),
		key:        aws.String(options.Key),
		downloader: options.Downloader,
		uploader:   options.Uploader,
	}
}

func (s S3StatsService) Get() (map[string]int, error) {
	data, err := s.downloader.Download()
	if err != nil {
		return nil, errors.Wrap(err, "error downloading")
	}

	var stats map[string]int
	err = json.Unmarshal(data, &stats)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling")
	}

	return stats, nil
}

func (s S3StatsService) Increment(member string) error {
	stats, err := s.Get()
	if err != nil {
		return errors.Wrap(err, "error getting stats")
	}

	stats[member] += 1

	data, err := json.Marshal(stats)
	if err != nil {
		return errors.Wrap(err, "error marshalling")
	}

	err = s.uploader.Upload(data)
	if err != nil {
		return errors.Wrap(err, "error uploading")
	}

	return nil
}
