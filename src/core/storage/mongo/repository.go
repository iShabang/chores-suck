package mongo

import (
	"chores-suck/core"
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

// NewStorage creates and returns a new storage object
func NewStorage(c *mongo.Client, ttl int32) *Storage {
	s := &Storage{
		cl: c,
	}
	s.setTTL("sessions", ttl)
	return s
}

// GetUserByName fetches a user from the database by unique username
func (s *Storage) GetUserByName(name string) (core.User, error) {
	filter := bson.M{"username": name}
	var u core.User
	err := s.findOne(&filter, u, "users")
	return u, err
}

// GetSession fetches a session frm the database by session id
func (s *Storage) GetSession(ses *sessions.Session) error {
	filter := bson.M{"ID": ses.ID}
	e := s.findOne(&filter, ses, "sessions")
	return e
}

// DeleteSession deletes a session from the database by ID
func (s *Storage) DeleteSession(ses *sessions.Session) error {
	filter := bson.M{"ID": ses.ID}
	e := s.deleteOne(&filter, "sessions")
	return e
}

// UpsertSession updates an existing session or inserts a new one.
func (s *Storage) UpsertSession(ses *sessions.Session) error {
	query := bson.M{"ID": ses.ID}
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

//func (s *Storage) findOne(filter *bson.M, object interface{}, coll string) (*mongo.SingleResult, error) {
func (s *Storage) findOne(filter *bson.M, object interface{}, coll string) error {
	collection := s.cl.Database("chores-suck").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res := collection.FindOne(ctx, filter)
	e := res.Decode(object)
	if e != nil && e != mongo.ErrNoDocuments {
		return e
	}
	return nil
}

func (s *Storage) deleteOne(filter *bson.M, coll string) error {
	var err error
	collection := s.cl.Database("chores-suck").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := collection.DeleteOne(ctx, filter)
	if err == nil && res.DeletedCount < 1 {
		err = errors.New("No items deleted")
	}
	return err
}

func (s *Storage) setTTL(collname string, expireTime int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	coll := s.cl.Database("chores-suck").Collection(collname)
	iv := coll.Indexes()
	opts := options.Index()
	opts.SetExpireAfterSeconds(expireTime)
	im := mongo.IndexModel{
		Keys:    bson.D{{Key: "MaxAge", Value: 1}},
		Options: opts,
	}
	_, err := iv.CreateOne(ctx, im)
	return err
}
