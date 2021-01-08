package mongo

import (
	"chores-suck/types"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Storage defines properties of a storage object
type Storage struct {
	cl *mongo.Client
}

// GetUserByName fetches a user from the database by unique username
func (s *Storage) GetUserByName(name string) (types.User, error) {
	filter := bson.D{{Key: "username", Value: name}}
	var u types.User
	res, err := s.findOne(&filter, "users")
	if err == nil {
		err = res.Decode(&u)
	}
	return u, err
}

func (s *Storage) findOne(filter *bson.D, coll string) (*mongo.SingleResult, error) {
	collection := s.cl.Database("chores-suck").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res := collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	return res, nil
}
