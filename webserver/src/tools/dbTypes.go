package tools

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

var ErrNotFound = errors.New("entry not found")

type UserSmall struct {
	Id string `json:"user_id"`
}

type UserLarge struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Username  string
	Attempts  uint8
}

type UserRecv struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	Username  string
	Attempts  int32
}

type Group struct {
	Id    string `json:"_id"`
	Admin string `json:"admin"`
	Name  string `json:"name"`
	Users []UserSmall
}

type Chore struct {
	Id      string `json:"_id"`
	Name    string `json:"name"`
	Time    uint   `json:"time"`
	UserId  string `json:"user_id"`
	GroupId string `json:"group_id"`
}

func (c *Chore) BsonD() bson.D {
	filter := bson.D{
		{Key: "name", Value: c.Name},
		{Key: "time", Value: c.Time},
		{Key: "user_id", Value: c.UserId},
		{Key: "group_id", Value: c.GroupId},
	}
	return filter
}

func (u *UserLarge) BsonD() bson.D {
	filter := bson.D{
		{Key: "first_name", Value: u.FirstName},
		{Key: "last_name", Value: u.LastName},
		{Key: "email", Value: u.Email},
		{Key: "password", Value: u.Password},
		{Key: "username", Value: u.Username},
		{Key: "attempts", Value: u.Attempts},
	}
	return filter
}

func (g *Group) BsonD() bson.D {
	filter := bson.D{
		{Key: "admin", Value: g.Admin},
		{Key: "name", Value: g.Name},
	}
	return filter
}
