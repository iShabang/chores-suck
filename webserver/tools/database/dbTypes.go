package db

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

/********************************************************
ERRORS
********************************************************/
var ErrNotFound = errors.New("entry not found")

/********************************************************
INTERFACE TYPES
********************************************************/
type DbType interface {
	BsonD() bson.D
}

/********************************************************
DATABASE OBJECT TYPES
********************************************************/
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

type Session struct {
	SessionId  string `bson:"sid"`
	UserId     string `bson:"uid"`
	ExpireTime int64  `bson:"exp"`
}

/********************************************************
BSON CONVERSION METHODS
********************************************************/
func (c Chore) BsonD() bson.D {
	filter := bson.D{
		{Key: "name", Value: c.Name},
		{Key: "time", Value: c.Time},
		{Key: "uid", Value: c.UserId},
		{Key: "gid", Value: c.GroupId},
	}
	return filter
}

func (u User) BsonD() bson.D {
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

func (g Group) BsonD() bson.D {
	filter := bson.D{
		{Key: "admin", Value: g.Admin},
		{Key: "name", Value: g.Name},
	}
	return filter
}

func (s Session) BsonD() bson.D {
	filter := bson.D{
		{Key: "sid", Value: s.SessionId},
		{Key: "uid", Value: s.UserId},
		{Key: "exp", Value: s.ExpireTime},
	}
	return filter
}
