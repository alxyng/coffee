package stats

import (
	"bytes"
	"errors"
	"testing"
)

type TestDownloader struct {
	data []byte
	err  error
}

func NewTestDownloader(data []byte, err error) *TestDownloader {
	return &TestDownloader{
		data: data,
		err:  err,
	}
}

func (u TestDownloader) Download() ([]byte, error) {
	return u.data, u.err
}

type TestUploader struct {
	Data []byte
	err  error
}

func NewTestUploader(err error) *TestUploader {
	return &TestUploader{
		err: err,
	}
}

func (u *TestUploader) Upload(data []byte) error {
	u.Data = data
	return u.err
}

func TestS3StatsServiceGetWhenDownloaderErrors(t *testing.T) {
	downloader := NewTestDownloader(nil, errors.New("some error"))
	service := createService(downloader, nil)

	stats, err := service.Get()

	if err == nil {
		t.Errorf("expected non nil err")
	}

	if stats != nil {
		t.Errorf("expected nil stats, got %v", stats)
	}
}

func TestS3StatsServiceGetWhenDownloaderReturnsEmptyData(t *testing.T) {
	data := []byte{}
	downloader := NewTestDownloader(data, nil)
	service := createService(downloader, nil)

	stats, err := service.Get()

	if err == nil {
		t.Errorf("expected non nil err")
	}

	if stats != nil {
		t.Errorf("expected nil stats, got %v", stats)
	}
}

func TestS3StatsServiceGetWhenDownloaderReturnsDataThatCannotBeMarshalled(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x04}
	downloader := NewTestDownloader(data, nil)
	service := createService(downloader, nil)
	stats, err := service.Get()

	if err == nil {
		t.Errorf("expected non nil err")
	}

	if stats != nil {
		t.Errorf("expected nil stats, got %v", stats)
	}
}

func TestS3StatsServiceGetWhenDownloaderReturnsDataForAnJsonObject(t *testing.T) {
	data := []byte("{}")
	downloader := NewTestDownloader(data, nil)
	service := createService(downloader, nil)

	stats, err := service.Get()

	if err != nil {
		t.Errorf("expected nil err, got %v", err)
	}

	if stats == nil {
		t.Errorf("expected non nil stats")
	}

	actualLength := len(stats)
	expectedLength := 0
	if actualLength != expectedLength {
		t.Errorf("incorrect length stats, got %v, want %v",
			actualLength, expectedLength)
	}
}

func TestS3StatsServiceGetWhenDownloaderReturnsDataThatCanBeMarshalled(t *testing.T) {
	data := []byte("{\"foo\":42}")
	downloader := NewTestDownloader(data, nil)
	service := createService(downloader, nil)

	stats, err := service.Get()

	if err != nil {
		t.Errorf("expected nil err, got %v", err)
	}

	if stats == nil {
		t.Errorf("expected non nil stats")
	}

	actualLength := len(stats)
	expectedLength := 1
	if actualLength != expectedLength {
		t.Errorf("incorrect length stats, got %v, want %v",
			actualLength, expectedLength)
	}

	actualValue := stats["foo"]
	expectedValue := 42
	if actualValue != expectedValue {
		t.Errorf("incorrect stats value, got %v, want %v",
			actualValue, expectedValue)
	}
}

func TestS3StatsServiceIncrementWhenDownloaderErrors(t *testing.T) {
	downloader := NewTestDownloader(nil, errors.New("some error"))
	service := createService(downloader, nil)

	err := service.Increment("")

	if err == nil {
		t.Errorf("expected non nil err")
	}
}

func TestS3StatsServiceIncrementWhenUploaderErrors(t *testing.T) {
	data := []byte("{\"foo\":42}")
	downloader := NewTestDownloader(data, nil)
	uploader := NewTestUploader(errors.New("some error"))
	service := createService(downloader, uploader)

	err := service.Increment("")

	if err == nil {
		t.Errorf("expected non nil err")
	}
}

func TestS3StatsServiceIncrementExistingMemberWhenUploaderSucceeds(t *testing.T) {
	data := []byte("{\"foo\":42}")
	downloader := NewTestDownloader(data, nil)
	uploader := NewTestUploader(nil)
	service := createService(downloader, uploader)

	err := service.Increment("foo")

	if err != nil {
		t.Errorf("expected nil err")
	}

	actualData := uploader.Data
	expectedData := []byte("{\"foo\":43}")
	if bytes.Compare(actualData, expectedData) != 0 {
		t.Errorf("incorrect data, got %s, want %s", actualData, expectedData)
	}
}

func TestS3StatsServiceIncrementNewMemberWhenUploaderSucceeds(t *testing.T) {
	data := []byte("{\"foo\":42}")
	downloader := NewTestDownloader(data, nil)
	uploader := NewTestUploader(nil)
	service := createService(downloader, uploader)

	err := service.Increment("bar")

	if err != nil {
		t.Errorf("expected nil err")
	}

	actualData := uploader.Data
	expectedData := []byte("{\"bar\":1,\"foo\":42}")
	if bytes.Compare(actualData, expectedData) != 0 {
		t.Errorf("incorrect data, got %s, want %s", actualData, expectedData)
	}
}

func createService(downloader Downloader, uploader Uploader) *S3StatsService {
	return NewS3StatsService(S3StatsOptions{
		Bucket:     "test-bucket",
		Key:        "results.json",
		Downloader: downloader,
		Uploader:   uploader,
	})
}
