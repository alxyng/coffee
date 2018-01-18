package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/boltdb/bolt"
)

type StatsService interface {
	Get() (map[string]int, error)
	Increment(member string) error
}

type MemoryStatsService struct {
	table map[string]int
}

func NewMemoryStatsService() MemoryStatsService {
	return MemoryStatsService{
		table: make(map[string]int),
	}
}

func (s MemoryStatsService) Get() (map[string]int, error) {
	return s.table, nil
}

func (s MemoryStatsService) Increment(member string) error {
	s.table[member] += 1
	return nil
}

type DiskStatsService struct {
	db *bolt.DB
}

func NewDiskStatsService(db *bolt.DB) (*DiskStatsService, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("members"))
		return err
	})

	if err != nil {
		return nil, err
	}

	return &DiskStatsService{
		db: db,
	}, nil
}

func (s DiskStatsService) Get() (map[string]int, error) {
	table := make(map[string]int)

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("members"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			n, err := strconv.ParseInt(string(v), 0, 0)
			if err != nil {
				return fmt.Errorf("could not parse int: %v", err)
			}
			table[string(k)] = int(n)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return table, nil
}

func (s DiskStatsService) Increment(member string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("members"))
		v := b.Get([]byte(member))
		var n int64
		if v == nil {
			n = 1
		} else {
			val, err := strconv.ParseInt(string(v), 0, 0)
			if err != nil {
				return fmt.Errorf("could not parse int: %v", err)
			}
			val++
			n = val
		}
		v = []byte(strconv.FormatInt(n, 10))
		return b.Put([]byte(member), v)
	})
}

type S3StatsOptions struct {
	Bucket     string
	Key        string
	Downloader *s3manager.Downloader
	Uploader   *s3manager.Uploader
}

type S3StatsService struct {
	bucket     *string
	key        *string
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
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
	buffer := &aws.WriteAtBuffer{}

	input := &s3.GetObjectInput{
		Bucket: s.bucket,
		Key:    s.key,
	}

	_, err := s.downloader.Download(buffer, input)
	if err != nil {
		return nil, err
	}

	var table map[string]int
	err = json.Unmarshal(buffer.Bytes(), &table)
	if err != nil {
		return nil, err
	}

	return table, nil
}

func (s S3StatsService) Increment(member string) error {
	table, err := s.Get()
	if err != nil {
		return err
	}

	table[member] += 1

	data, err := json.Marshal(table)
	if err != nil {
		return err
	}

	reader := bytes.NewBuffer(data)

	input := &s3manager.UploadInput{
		Bucket: s.bucket,
		Key:    s.key,
		Body:   reader,
	}

	_, err = s.uploader.Upload(input)
	if err != nil {
		return err
	}

	return nil
}
