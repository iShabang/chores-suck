package tools

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type UserSmall struct {
	userId string
	name   string
}

type UserLarge struct {
	firstName string
	lastName  string
	email     string
	password  string
	username  string
}

type Group struct {
	id    string
	admin string
	name  string
	users []UserSmall
}

type Chore struct {
	Id      string `json:"_id"`
	Name    string `json:"name"`
	Time    uint   `json:"time"`
	UserId  string `json:"user_id"`
	GroupId string `json:"group_id"`
}

// get chores for a group
// db.chores.find({group_id: ""})
func (c *Connection) GetGroupChores(id string) ([]*Chore, error) {
	fmt.Printf("looking for id %v\n", id)
	collection := c.client.Database("fairmate").Collection("chores")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	filter := bson.M{"group_id": id}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer cur.Close(ctx)
	var chores []*Chore
	fmt.Print("starting loop\n")
	for cur.Next(ctx) {
		var chore Chore
		err := cur.Decode(&chore)
		if err != nil {
			fmt.Print(err)
			return nil, err
		}
		chores = append(chores, &chore)
		fmt.Printf("got item %v", chore.Name)
	}
	return chores, nil
}

// get chores for a user
// db.chores.find({user_id: ""})

// add chores to a group
/*
db.chores.insertMany([
{},
{},
{},
])
*/

// add user
/*
db.users.insert({
    first_name: "Shannon",
    last_name: "Curtin",
    email: "curtin.shannon@gmail.com",
    password: "1234",
    username: "iShaBaNg"
})
*/

// add user to group
/*
db.groups.update({_id: ""},
{
	$push:{
		users: {
			user_id: ""
		}
	}
})
*/

// add group
/*
db.groups.insert({
    admin: "5df2ae6a0aba119155a01b8c"
    name: "UNLV House",
    users: [
        {user_id:"5df2ae6a0aba119155a01b8c"}
    ]
})
*/
