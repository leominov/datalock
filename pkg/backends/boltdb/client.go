package boltdb

import (
	"encoding/json"
	"errors"
	"path"

	"github.com/boltdb/bolt"
)

type Client struct {
	bucket []byte
	db     *bolt.DB
}

func NewClient(directory, bucket string) (*Client, error) {
	client := &Client{
		bucket: []byte(bucket),
	}
	db, err := bolt.Open(path.Join(directory, "datalock.db"), 0600, nil)
	if err != nil {
		return nil, err
	}
	client.db = db
	err = client.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(client.bucket); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) GetValue(id string, body interface{}) error {
	return c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucket)
		v := b.Get([]byte(id))
		if len(v) == 0 {
			return errors.New("Meta not found")
		}
		if err := json.Unmarshal(v, &body); err != nil {
			return err
		}
		return nil
	})
}

func (c *Client) SetValue(id string, body interface{}) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucket)
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		return b.Put([]byte(id), encoded)
	})
}

func (c *Client) Close() error {
	return c.db.Close()
}
