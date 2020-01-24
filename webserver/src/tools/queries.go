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
	filter := bson.D{{Key: "group_id", Value: id}}
	return c.getChores(&filter)

}

func (c *Connection) getChores(f *bson.D) ([]*Chore, error) {
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
	filter := bson.D{{Key: "user_id", Value: id}}
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
	filter := ch.BsonD()
	return c.insert(&filter, "chores")
}

func (c *Connection) AddManyChores(ch []Chore) error {
	temp := make([]bson.D, len(ch))
	for i, v := range ch {
		temp[i] = v.BsonD()
	}
	return c.insertMany(temp, "chores")
}

func (c *Connection) insert(f *bson.D, coll string) (string, error) {
	collection := c.client.Database("fairmate").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, f)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	id := (res.InsertedID).(primitive.ObjectID)
	return id.String(), nil
}

func (c *Connection) insertMany(f []bson.D, coll string) error {
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
func (c *Connection) AddUser(u *User) (string, error) {
	filter := u.BsonD()
	collection := c.client.Database("fairmate").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, filter)
	id := (res.InsertedID).(primitive.ObjectID)
	return id.String(), err
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
	groupFilter := bson.D{{Key: "_id", Value: gId}}
	updateFilter := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "users", Value: uId}}}}
	return c.update(&groupFilter, &updateFilter, "groups")
}

func (c *Connection) update(gf *bson.D, of *bson.D, coll string) error {
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
	filter := g.BsonD()
	return c.insert(&filter, "groups")
}

func (c *Connection) GetUser(username string) (*User, error) {
	filter := bson.D{{Key: "username", Value: username}}
	return c.getUser(&filter)
}

func (c *Connection) UpdateUserAttempts(u string, attempts int32) error {
	userFilter := bson.D{{Key: "username", Value: u}}
	updateFilter := bson.D{{Key: "$set", Value: bson.D{{Key: "attempts", Value: attempts}}}}
	return c.update(&userFilter, &updateFilter, "users")
}

func (c *Connection) getUser(filter *bson.D) (*User, error) {
	collection := c.client.Database("fairmate").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var u User
	err := collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("found document %v\n", u)
	/*
		id := (b["_id"]).(primitive.ObjectID).String()
		u.Id = id
		u.FirstName = (b["firstname"]).(string)
		u.LastName = (b["lastname"]).(string)
		u.Email = (b["email"]).(string)
		u.Password = (b["password"]).(string)
		u.Username = (b["username"]).(string)
		u.Attempts = (b["attempts"]).(int32)
	*/
	return &u, nil
}

func (c *Connection) AddSession(u *User, id string, t time.Time) (string, error) {
	exp := fmt.Sprintf("%v", t.Unix())
	filter := bson.D{{Key: "sid", Value: id}, {Key: "uid", Value: u.Id}, {Key: "exp", Value: exp}}
	return c.insert(&filter, "session")
}
