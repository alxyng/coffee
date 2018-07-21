package stats

import (
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
)

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
