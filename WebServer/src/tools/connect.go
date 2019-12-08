package tools

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Connection struct {
	logger *log.Logger
	url    string
}

func NewConnection() Connection {
	return Connection{}
}

func (c *Connection) Connect(u string) bool {
	ret := false
	client, err := mongo.NewClient(options.Client().ApplyURI(u))
	if err == nil {
		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		err = client.Connect(ctx)
		if err == nil {
			ret = true
		} else {
			c.logger.Print(err)
		}
	} else {
		c.logger.Print(err)
	}
	return ret
}
