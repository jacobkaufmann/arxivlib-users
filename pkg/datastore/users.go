package datastore

import (
	"context"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *usersStore) Authenticate(username, passwd string) (*arxivlib.User, error) {
	coll := s.db.Collection("users")
	var user *arxivlib.User

	filter := bson.M{"username": username, "passwd": passwd}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *usersStore) List(opt *arxivlib.UserListOptions) ([]*arxivlib.User, error) {
	if opt == nil {
		opt = &arxivlib.UserListOptions{}
	}

	coll := s.db.Collection("users")
	var users []*arxivlib.User

	filter := bson.M{
		"username": primitive.Regex{Pattern: opt.Username, Options: "i"},
	}
	projection := bson.M{
		"_id":      1,
		"username": 1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := coll.Find(
		ctx,
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

func (s *usersStore) Create(user *arxivlib.User) (created bool, err error) {
	coll := s.db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := coll.InsertOne(ctx, &user)
	if err != nil {
		return
	}

	id := result.InsertedID
	if id != nil {
		created = true
	}

	return
}

func (s *usersStore) Delete(id primitive.ObjectID) (deleted bool, err error) {
	coll := s.db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return
	}

	if result.DeletedCount > 0 {
		deleted = true
	}

	return
}
