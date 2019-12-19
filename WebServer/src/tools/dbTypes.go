package tools

import (
	"go.mongodb.org/mongo-driver/bson"
)

type UserSmall struct {
	Id string `json:"user_id"`
}

type UserLarge struct {
	Id        string `json:"_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
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

func (c *Chore) BsonM() bson.M {
	filter := bson.M{
		"name":     c.Name,
		"time":     c.Time,
		"user_id":  c.UserId,
		"group_id": c.GroupId,
	}
	return filter
}

func (u *UserLarge) BsonM() bson.M {
	filter := bson.M{
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"email":      u.Email,
		"password":   u.Password,
		"username":   u.Username,
	}
	return filter
}

func (g *Group) BsonM() bson.M {
	filter := bson.M{
		"admin": g.Admin,
		"name":  g.Name,
	}
	return filter
}
