// db
package utils

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
)

type DBArch struct {
	Title            string `json:"title,omitempty"`
	Url              string `json:"url,omitempty"`
	CreateTime       int    `json:"create_time,omitempty"`
	UpdateTime       int    `json:"update_time,omitempty"`
	IsSubmitDownload bool   `json:"is_submit_download,omitempty"`
}

func DBInit() {

}

func DBUpdate(t *DBArch, key string) error {
	db, err := bolt.Open("data.db", 0600, nil)
	defer db.Close()
	if err != nil {
		return err
	}
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, err := tx.CreateBucketIfNotExists([]byte("hardseed"))
	if err != nil {
		return err
	}
	encoded, err := json.Marshal(t)
	if err != nil {
		return err
	}
	err = b.Put([]byte(key), []byte(encoded))
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func DBGet(key string) (*DBArch, error) {
	db, err := bolt.Open("data.db", 0600, nil)
	defer db.Close()
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b, err := tx.CreateBucketIfNotExists([]byte("hardseed"))
	if err != nil {
		return nil, err
	}
	v := b.Get([]byte(key))
	if v == nil {
		return nil, nil
	}
	data := new(DBArch)
	if err := json.Unmarshal(v, data); err != nil {
		log.Println("JSON Unmarshal error:", err)
	}
	return data, nil
}
