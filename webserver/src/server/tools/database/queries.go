package db

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

/********************************************************
EXPORTED API
********************************************************/

///////////////////// INSERT ////////////////////////////
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

func (c *Connection) AddUser(u *User) (string, error) {
	filter := u.BsonD()
	collection := c.client.Database("fairmate").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, filter)
	id := (res.InsertedID).(primitive.ObjectID)
	return id.String(), err
}

func (c *Connection) AddUserToGroup(uId string, gId string) error {
	groupFilter := bson.D{{Key: "_id", Value: gId}}
	updateFilter := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "users", Value: uId}}}}
	return c.update(&groupFilter, &updateFilter, "groups")
}

func (c *Connection) AddGroup(g *Group) (string, error) {
	filter := g.BsonD()
	return c.insert(&filter, "groups")
}

func (c *Connection) AddSession(s *Session) (string, error) {
	filter := s.BsonD()
	return c.insert(&filter, "session")
}

///////////////////// UPDATE ////////////////////////////
func (c *Connection) UpdateUserAttempts(u string, attempts int32) error {
	userFilter := bson.D{{Key: "username", Value: u}}
	updateFilter := bson.D{{Key: "$set", Value: bson.D{{Key: "attempts", Value: attempts}}}}
	return c.update(&userFilter, &updateFilter, "users")
}

///////////////////// FIND ////////////////////////////
func (c *Connection) GetGroupChores(id string) ([]Chore, error) {
	result := true
	filter := bson.D{{Key: "gid", Value: id}}

	objs, err := c.findMany(&filter, "chores")
	result = (err == nil)
	chs := make([]Chore, len(objs))
	if result {
		for i, v := range objs {
			bson.Unmarshal(v, chs[i])
		}
	}
	return chs, err
}

func (c *Connection) GetUserChores(id string) ([]*Chore, error) {
	filter := bson.D{{Key: "user_id", Value: id}}
	return c.getChores(&filter)
}

func (c *Connection) GetUser(username string) (*User, error) {
	filter := bson.D{{Key: "username", Value: username}}
	var u User
	err := c.findOne(&filter, "users", u)
	return &u, err
}

func (c *Connection) FindSession(sid string) (*Session, error) {
	filter := bson.D{{Key: "sid", Value: sid}}
	var sess Session
	err := c.findOne(&filter, "session", sess)
	return &sess, err
}

///////////////////// DELETE ////////////////////////////
func (c *Connection) DeleteSessions(userId string) error {
	filter := bson.D{{Key: "uid", Value: userId}}
	return c.deleteMany(&filter, "session")
}

/********************************************************
INTERNAL METHODS
********************************************************/

/********************************************************
INSERTIONS
********************************************************/
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

/********************************************************
UPDATE
********************************************************/
func (c *Connection) update(gf *bson.D, of *bson.D, coll string) error {
	var err error
	collection := c.client.Database("fairmate").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := collection.UpdateOne(ctx, gf, of)
	if err == nil && res.ModifiedCount < 1 {
		err = errors.New("no objects modified")
	}
	return err
}

/********************************************************
FIND
********************************************************/
func (c *Connection) findOne(filter *bson.D, coll string, obj DbType) error {
	collection := c.client.Database("fairmate").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(obj)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("found document %v\n", obj)
	return nil
}

func (c *Connection) findMany(f *bson.D, coll string) ([]bson.Raw, error) {
	collection := c.client.Database("fairmate").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, f)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer cur.Close(ctx)
	var temp []bson.Raw
	fmt.Print("starting loop\n")
	for cur.Next(ctx) {
		temp = append(temp, cur.Current)
	}
	return temp, nil
}

// Not Used
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

// Not used
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
	return &u, nil
}

/********************************************************
DELETE
********************************************************/
func (c *Connection) deleteMany(f *bson.D, coll string) error {
	var err error
	collection := c.client.Database("fairmate").Collection(coll)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := collection.DeleteMany(ctx, f)
	if err == nil && res.DeletedCount < 1 {
		err = errors.New("No items deleted")
	}
	return err
}
