package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Work with a specific MongoDB collection
type MongoClientCollection struct {
	col string
	d   *MongoClientDatabase
}

func (c *MongoClientCollection) Database() *MongoClientDatabase {
	return c.d
}

func NewMongoClientCollection(d *MongoClientDatabase, collection string) *MongoClientCollection {
	return &MongoClientCollection{d: d, col: collection}
}

func (c *MongoClientCollection) InsertOne(obj interface{}) (*mongo.InsertOneResult, error) {
	return c.d.mc.InsertOne(c.d.db, c.col, obj)
}

func (c *MongoClientCollection) InsertMany(objs []interface{}) (*mongo.InsertManyResult, error) {
	return c.d.mc.InsertMany(c.d.db, c.col, objs)
}

// obj MUST be a pointer
func (c *MongoClientCollection) FindID(id string, obj interface{}) error {
	return c.d.mc.FindID(c.d.db, c.col, id, obj)
}

// obj MUST be a pointer
// find by the primitive ObjectID
func (c *MongoClientCollection) FindObjectID(id primitive.ObjectID, obj interface{}) error {
	return c.d.mc.FindObjectID(c.d.db, c.col, id, obj)
}

// obj MUST be a pointer
func (c *MongoClientCollection) FindOne(filter, obj interface{}) error {
	return c.d.mc.FindOne(c.d.db, c.col, filter, obj)
}

func (c *MongoClientCollection) Find(filter, slice interface{}) error {
	return c.d.mc.Find(c.d.db, c.col, filter, slice)
}

func (c *MongoClientCollection) UpdateOne(filter, update interface{}) error {
	return c.d.mc.UpdateOne(c.d.db, c.col, filter, update)
}

func (c *MongoClientCollection) UpdateMany(filter, update interface{}) error {
	return c.d.mc.UpdateOne(c.d.db, c.col, filter, update)
}

func (c *MongoClientCollection) DeleteOne(filter interface{}) error {
	return c.d.mc.DeleteOne(c.d.db, c.col, filter)
}

func (c *MongoClientCollection) DeleteMany(filter interface{}) error {
	return c.d.mc.DeleteMany(c.d.db, c.col, filter)
}
