package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const requestTimeout = 5 * time.Second

// Work with MongoDB in its entierty
type MongoClient struct {
	ctx context.Context
	uri string
	mc  *mongo.Client
}

func NewMongoClient(parent context.Context, uri string) *MongoClient {
	return &MongoClient{
		ctx: parent,
		uri: uri,
		mc:  &mongo.Client{},
	}
}

// return a MongoClientDatabase instance
func (c *MongoClient) Database(db string) *MongoClientDatabase {
	return NewMongoClientDatabase(c, db)
}

func (c *MongoClient) ListDatabaseNames(filter interface{}) ([]string, error) {
	// make a timeout to request database names
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()
	return c.mc.ListDatabaseNames(ctx, filter)
}

func (c *MongoClient) ListCollectionNames(db string, filter interface{}) ([]string, error) {
	// make a timeout to request database names
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()
	return c.mc.Database(db).ListCollectionNames(ctx, filter)
}

// db: where to create the new collection
// collection: the name of the new collection
func (c *MongoClient) CreateCollection(db, collection string) error {
	// make a timeout to request collection creation
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()
	return c.mc.Database(db).CreateCollection(ctx, collection)
}

func (c *MongoClient) DropCollection(db, collection string) error {
	col := c.mc.Database(db).Collection(collection)
	// make a timeout to request collection drop
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()
	return col.Drop(ctx)
}

func (c *MongoClient) DropDatabase(db, collection string) error {
	database := c.mc.Database(db)
	// make a timeout to request collection drop
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()
	return database.Drop(ctx)
}

func (c *MongoClient) Connect() error {
	// make a timeout to connect and ping mongodb
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()

	// connect to the mongodb server
	mc, err := mongo.Connect(ctx, options.Client().ApplyURI(c.uri))
	if err != nil {
		return fmt.Errorf("failed to connect mongodb: %w", err)
	}
	c.mc = mc

	// test if the connection has been successfully established
	return c.mc.Ping(ctx, nil)
}

func (c *MongoClient) DropAndCreateIndex(db, collection, column string, nb uint, unique bool) (string, error) {
	col := c.mc.Database(db).Collection(collection)
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()
	// drop if already exists
	col.Indexes().DropOne(ctx, fmt.Sprintf("%s_%d", column, nb))
	return col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{column: nb},
		Options: options.Index().SetUnique(unique),
	})
}

func (c *MongoClient) CreateIndex(db, collection, column string, nb uint, unique bool) (string, error) {
	col := c.mc.Database(db).Collection(collection)
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()
	return col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{column: nb},
		Options: options.Index().SetUnique(unique),
	})
}

func (c *MongoClient) InsertOne(db, collection string, obj interface{}) (*mongo.InsertOneResult, error) {
	col := c.mc.Database(db).Collection(collection)
	// find one with a timeout
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()
	return col.InsertOne(ctx, obj)
}

func (c *MongoClient) InsertMany(db, collection string, objs []interface{}) (*mongo.InsertManyResult, error) {
	col := c.mc.Database(db).Collection(collection)
	// find one with a timeout
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()
	return col.InsertMany(ctx, objs)
}

// obj MUST be a pointer
// obj is an instance of the object type the request will be decoded in
func (c *MongoClient) FindID(db, collection, id string, obj interface{}) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return c.FindOne(db, collection, bson.M{"_id": objectId}, obj)
}

// obj MUST be a pointer
// obj is an instance of the object type the request will be decoded in
// find directly by a primitive ObjectID and not a string converted to it
func (c *MongoClient) FindObjectID(db, collection string, objectId primitive.ObjectID, obj interface{}) error {
	return c.FindOne(db, collection, bson.M{"_id": objectId}, obj)
}

// obj MUST be a pointer
// obj is an instance of the object type the request will be decoded in
func (c *MongoClient) FindOne(db, collection string, filter, obj interface{}) error {
	// get the db in which our collection is
	// then get the collection
	col := c.mc.Database(db).Collection(collection)
	// find one with a timeout
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()
	result := col.FindOne(ctx, filter)
	if result.Err() != nil {
		return result.Err()
	}
	return result.Decode(obj)
}

// obj MUST be a pointer
// obj is an instance of the object type the request will be decoded in
//
// returns a slice of obj type instances
func (c *MongoClient) Find(db, collection string, filter, slice interface{}) error {
	// get the db in which our collection is
	// then get the collection
	col := c.mc.Database(db).Collection(collection)
	// find one with a timeout
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()
	cur, err := col.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	return cur.All(ctx, slice)
}

// updatedFields must be a bson.D object of updated key, value
// e.g. bson.D{{"name", "new name"}, ...}
func (c *MongoClient) UpdateOne(db, collection string, filter, updatedFields interface{}) error {
	// get the db in which our collection is
	// then get the collection
	col := c.mc.Database(db).Collection(collection)
	// find one with a timeout
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()

	_, err := col.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updatedFields}})
	return err
}

// updatedFields must be a bson.D object of updated key, value
// e.g. bson.D{{"name", "new name"}, ...}
func (c *MongoClient) UpdateMany(db, collection string, filter, updatedFields interface{}) error {
	// get the db in which our collection is
	// then get the collection
	col := c.mc.Database(db).Collection(collection)
	// find one with a timeout
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()

	_, err := col.UpdateMany(ctx, filter, bson.D{{Key: "$set", Value: updatedFields}})
	return err
}

func (c *MongoClient) DeleteOne(db, collection string, filter interface{}) error {
	// get the db in which our collection is
	// then get the collection
	col := c.mc.Database(db).Collection(collection)
	// find one with a timeout
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()

	_, err := col.DeleteOne(ctx, filter)
	return err
}

func (c *MongoClient) DeleteMany(db, collection string, filter interface{}) error {
	// get the db in which our collection is
	// then get the collection
	col := c.mc.Database(db).Collection(collection)
	// find one with a timeout
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()

	_, err := col.DeleteMany(ctx, filter)
	return err
}

func (c *MongoClient) Disconnect() error {
	// disconnect from the mongodb server
	ctx, cancel := context.WithTimeout(c.ctx, requestTimeout)
	defer cancel()
	return c.mc.Disconnect(ctx)
}
