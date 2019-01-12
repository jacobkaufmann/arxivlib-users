package datastore

import (
	arxivlib "github.com/jacobkaufmann/arxivlib-users"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// A Datastore accesses the datastore (MongoDB)
type Datastore struct {
	Users arxivlib.UsersService

	db *mongo.Database
}

// NewDatastore creates a new client for accessing the datastore
func NewDatastore(db *mongo.Database) *Datastore {
	if db == nil {
		db = DB
	}
	d := &Datastore{db: db}
	d.Users = &usersStore{d}

	return d
}
