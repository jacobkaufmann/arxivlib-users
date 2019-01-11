package datastore

import (
	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo/options"

	arxivlib "github.com/jacobkaufmann/arxivlib-users"
	"github.com/mongodb/mongo-go-driver/bson"
)

type usersStore struct {
	*Datastore
}

func (s *usersStore) Get(id primitive.ObjectID) (*arxivlib.User, error) {
	coll := s.db.Collection("users")
	var user *arxivlib.User

	filter := bson.M{"_id": id}

	err := coll.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *usersStore) Authenticate(username, passwd string) (*arxivlib.User, error) {
	coll := s.db.Collection("users")
	var user *arxivlib.User

	filter := bson.M{"username": username, "passwd": passwd}

	err := coll.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *usersStore) List(opt *arxivlib.UserListOptions) ([]*arxivlib.User, error) {
	if s.db == nil {
		s.db = DB
	}
	if opt == nil {
		opt = &arxivlib.UserListOptions{}
	}

	coll := s.db.Collection("users")
	var users []*arxivlib.User

	filter := bson.D{
		{"username", primitive.Regex{Pattern: opt.Username, Options: "i"}},
	}
	projection := bson.D{
		{"_id", 1},
		{"username", 1},
	}

	cursor, err := coll.Find(
		context.Background(),
		filter,
		options.Find().SetProjection(projection),
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		user := &arxivlib.User{}
		if err := cursor.Decode(user); err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *usersStore) Create(user *arxivlib.User) (bool, error) {
	coll := s.db.Collection("users")

	result, err := coll.InsertOne(
		context.Background(),
		&user,
	)
	if err != nil {
		return false, err
	}

	id := result.InsertedID
	log.Printf("User created with _id: %v\n", id)

	return true, nil
}
