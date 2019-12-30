package tools

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type Connection struct {
	logger *log.Logger
	url    string
	client *mongo.Client
}

func NewConnection() Connection {
	return Connection{}
}

func (c *Connection) Connect(u string) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(u))
	if err == nil {
		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		err = client.Connect(ctx)
		if err == nil {
			err = client.Ping(ctx, readpref.Primary())
			if err == nil {
				c.client = client
			}
		}
	}
	return err
}
