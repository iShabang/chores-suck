package tools

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

var ErrNotFound = errors.New("entry not found")

type User struct {
	Id        string `bson:"_id"`
	FirstName string `bson:"firstname"`
	LastName  string `bson:"lastname"`
	Email     string `bson:"email"`
	Password  string `bson:"password"`
	Username  string `bson:"username"`
	Attempts  int32  `bson:"attempts"`
}

type Group struct {
	Id    string   `bson:"_id"`
	Admin string   `bson:"admin"`
	Name  string   `bson:"name"`
	Users []string `bson:"users"`
}

type Chore struct {
	Id      string `bson:"_id"`
	Name    string `bson:"name"`
	Time    uint   `bson:"time"`
	UserId  string `bson:"userid"`
	GroupId string `bson:"groupid"`
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

func (u *User) BsonD() bson.D {
	filter := bson.D{
		{Key: "firstname", Value: u.FirstName},
		{Key: "lastname", Value: u.LastName},
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
