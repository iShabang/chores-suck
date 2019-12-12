package tools

import (
	"context"
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
	id     string
	name   string
	time   uint
	userId string
}

// get chores for a group
// db.chores.find({group_id: ""})
func (c *Connection) GetGroupChores(id string) ([]Chore, error) {
	collection := c.client.Database("fairmate").Collection("chores")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	filter := bson.M{"group_id": id}
	cur, err := collection.Find(ctx, filter)
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err == nil {

		}

	}
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
