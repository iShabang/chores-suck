package tools

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

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
	filter := ch.BsonM()
	return c.insert(&filter, "chores")
}

func (c *Connection) AddManyChores(ch []Chore) error {
	temp := make([]bson.M, len(ch))
	for i, v := range ch {
		temp[i] = v.BsonM()
	}
	return c.insertMany(temp, "chores")
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

func (c *Connection) insertMany(f []bson.M, coll string) error {
	ui := make([]interface{}, len(f))
	for i, v := range f {
		ui[i] = v
	}
	collection := c.client.Database("fairmate").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := collection.InsertMany(ctx, ui)
	return err
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
	filter := u.BsonM()
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
	return c.update(&groupFilter, &updateFilter, "groups")
}

func (c *Connection) update(gf *bson.M, of *bson.M, coll string) error {
	collection := c.client.Database("fairmate").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := collection.UpdateOne(ctx, gf, of)
	if err != nil {
		fmt.Print(err)
		return err
	}
	if res.ModifiedCount < 1 {
		return errors.New("no objects modified")
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

func (c *Connection) AddGroup(g *Group) (string, error) {
	filter := g.BsonM()
	return c.insert(&filter, "groups")
}

func (c *Connection) GetUser(username string) (*UserLarge, error) {
	filter := bson.M{"username": username}
	return c.getUser(&filter)
}

func (c *Connection) UpdateUserAttempts(u string, attempts uint8) error {
	userFilter := bson.M{"username": u}
	updateFilter := bson.M{"$set": bson.M{"attempts": attempts}}
	return c.update(&userFilter, &updateFilter, "users")
}

func (c *Connection) getUser(filter *bson.M) (*UserLarge, error) {
	collection := c.client.Database("fairmate").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, ErrNotFound
	}
	defer cur.Close(ctx)
	var u UserLarge
	cur.Next(ctx)
	err = cur.Decode(&u)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("got item %v\n", u.FirstName)
	return &u, nil
}

func (c *Connection) AddSession(u *UserLarge, id string, t time.Time) (string, error) {
	filter := bson.M{"sid": id, "uid": u.Id, "exp": fmt.Sprintf("%v", t.Unix)}
	return c.insert(&filter, "session")
}
