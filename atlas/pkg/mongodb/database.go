package mongodb

import "go.mongodb.org/mongo-driver/mongo"

// Work with a specific MongoDB database
type MongoClientDatabase struct {
	db string
	mc *MongoClient
}

func NewMongoClientDatabase(mc *MongoClient, db string) *MongoClientDatabase {
	return &MongoClientDatabase{mc: mc, db: db}
}

func (d *MongoClientDatabase) Collection(collection string) *MongoClientCollection {
	return NewMongoClientCollection(d, collection)
}

func (d *MongoClientDatabase) InsertOne(collection string, obj interface{}) (*mongo.InsertOneResult, error) {
	return d.mc.InsertOne(d.db, collection, obj)
}

func (d *MongoClientDatabase) InsertMany(collection string, objs []interface{}) (*mongo.InsertManyResult, error) {
	return d.mc.InsertMany(d.db, collection, objs)
}

// obj MUST be a pointer
func (d *MongoClientDatabase) FindID(collection, id string, obj interface{}) error {
	return d.mc.FindID(d.db, collection, id, obj)
}

// obj MUST be a pointer
func (d *MongoClientDatabase) FindOne(collection string, filter, obj interface{}) error {
	return d.mc.FindOne(d.db, collection, filter, obj)
}

func (d *MongoClientDatabase) Find(collection string, filter, slice interface{}) error {
	return d.mc.Find(d.db, collection, filter, slice)
}

func (d *MongoClientDatabase) UpdateOne(collection string, filter, update interface{}) error {
	return d.mc.UpdateOne(d.db, collection, filter, update)
}

func (d *MongoClientDatabase) UpdateMany(collection string, filter, update interface{}) error {
	return d.mc.UpdateOne(d.db, collection, filter, update)
}

func (d *MongoClientDatabase) DeleteOne(collection string, filter interface{}) error {
	return d.mc.DeleteOne(d.db, collection, filter)
}

func (d *MongoClientDatabase) DeleteMany(db, collection string, filter interface{}) error {
	return d.mc.DeleteMany(d.db, collection, filter)
}
