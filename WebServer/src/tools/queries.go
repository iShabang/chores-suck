package tools

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserSmall struct {
	Id   string
	name string
}

type UserLarge struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	Username  string
}

type Group struct {
	Id    string
	Admin string
	Name  string
	Users []UserSmall
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
	filter := bson.M{"group_id": id}
	return c.getChores(&filter)

}

func (c *Connection) getChores(f *bson.M) ([]*Chore, error) {
	collection := c.client.Database("fairmate").Collection("chores")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, f)
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
		fmt.Printf("got item %v\n", chore.Name)
	}
	return chores, nil
}

// get chores for a user
// db.chores.find({user_id: ""})
func (c *Connection) GetUserChores(id string) ([]*Chore, error) {
	filter := bson.M{"user_id": id}
	return c.getChores(&filter)
}

// add chores to a group
/*
db.chores.insertMany([
{},
{},
{},
])
*/
func (c *Connection) AddChore(ch *Chore) (string, error) {
	filter := bson.M{
		"name":     ch.Name,
		"time":     ch.Time,
		"user_id":  ch.UserId,
		"group_id": ch.GroupId,
	}
	return c.insert(&filter, "chores")
}

func (c *Connection) insert(f *bson.M, coll string) (string, error) {
	collection := c.client.Database("fairmate").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, f)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	id := (res.InsertedID).(primitive.ObjectID)
	fmt.Printf("added new chore with id: %v", id.String())
	return id.String(), nil
}

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
func (c *Connection) AddUser(u *UserLarge) (string, error) {
	filter := bson.M{
		"first_name": u.FirstName,
		"last_name":  u.LastName,
		"email":      u.Email,
		"password":   u.Password,
		"username":   u.Username,
	}
	return c.insert(&filter, "users")
}

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
func (c *Connection) AddUserToGroup(uId string, gId string) error {
	groupFilter := bson.M{"_id": gId}
	updateFilter := bson.M{"$addToSet": bson.M{"users": uId}}
	return c.updateGroup(&groupFilter, &updateFilter, "groups")
}

func (c *Connection) updateGroup(gf *bson.M, of *bson.M, coll string) error {
	collection := c.client.Database("fairmate").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := collection.UpdateOne(ctx, gf, of)
	if err != nil {
		fmt.Print(err)
		return err
	}
	if res.ModifiedCount < 1 {
		// create and return new error stating such
	}
	return nil
}

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
