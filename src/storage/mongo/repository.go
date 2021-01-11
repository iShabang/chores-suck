package mongo

import (
	"chores-suck/types"
	"context"
	"errors"
	"time"

	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// GetSession fetches a session frm the database by session id
func (s *Storage) GetSession(ses *sessions.Session) error {
	filter := bson.D{{Key: "sid", Value: ses.ID}}
	r, e := s.findOne(&filter, "sessions")
	if e == nil {
		e = r.Decode(ses)
	}
	return e
}

// DeleteSession deletes a session from the database by ID
func (s *Storage) DeleteSession(ses *sessions.Session) error {
	filter := bson.D{{Key: "sid", Value: ses.ID}}
	e := s.deleteOne(&filter, "sessions")
	return e
}

// UpsertSession updates an existing session or inserts a new one.
func (s *Storage) UpsertSession(ses *sessions.Session) error {
	query := bson.M{"sid": ses.ID}
	update, err := bson.Marshal(ses)
	var options options.UpdateOptions
	options.SetUpsert(true)
	s.upsert(&query, update, &options, "sessions")
	return err
}

func (s *Storage) upsert(query *bson.M, update []byte, options *options.UpdateOptions, coll string) error {
	collection := s.cl.Database("chores-suck").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := collection.UpdateOne(ctx, query, update, options)
	if err == nil && res.ModifiedCount < 1 {
		err = errors.New("no objects modified")
	}
	return err

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

func (s *Storage) deleteOne(f *bson.D, coll string) error {
	var err error
	collection := s.cl.Database("chores-suck").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := collection.DeleteOne(ctx, f)
	if err == nil && res.DeletedCount < 1 {
		err = errors.New("No items deleted")
	}
	return err
}
