package backends

import (
	"errors"

	"github.com/leominov/datalock/pkg/backends/boltdb"
)

type StoreClient interface {
	GetValue(id string, body interface{}) error
	SetValue(id string, body interface{}) error
	Close() error
}

func New(config Config) (StoreClient, error) {
	if config.Backend == "" {
		config.Backend = "boltdb"
	}
	switch config.Backend {
	case "boltdb":
		if config.Directory == "" {
			config.Directory = "./database"
		}
		if config.Bucket == "" {
			config.Bucket = "meta"
		}
		return boltdb.NewClient(config.Directory, config.Bucket)
	}
	return nil, errors.New("Invalid backend")
}
