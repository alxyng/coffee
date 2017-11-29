package main

import (
	"strconv"

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
			n, _ := strconv.ParseInt(string(v), 0, 0)
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
		n, _ := strconv.ParseInt(string(v), 0, 0)
		n++
		v = []byte(strconv.FormatInt(n, 10))
		return b.Put([]byte(member), v)
	})
}
