package dbfilekit

import "go.etcd.io/bbolt"

type DBFileKit interface {
	GetDB() (*bbolt.DB, error)
}
