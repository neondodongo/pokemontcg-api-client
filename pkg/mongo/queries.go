package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/dto"
)

type DataAccess interface {
	GetCollection() mongo.Collection
	Upsert(interface{}, string)
	GetCardById(id string) dto.Card
}

type MongoBongo struct {
	Client          *mongo.Client
	Database        string
	CardsCollection string
	SetsCollection  string
	UsersCollection string
}

func (db *MongoBongo) Upsert(t interface{}, c string) error {
	var filter bson.M

	//determine interface type
	switch t.(type) {
	case dto.Card:
		filter = bson.M{"id": t.(dto.Card).ID}
	case dto.Set:
		filter = bson.M{"code": t.(dto.Set).Code}
	case dto.User:
		filter = bson.M{"code": t.(dto.User).Username}
	}

	update := bson.M{"$set": t}
	userChosenCollection := db.GetCollection(c)
	r, err := userChosenCollection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("Failed to insert one to %s collection [ %v ]", userChosenCollection.Name(), err)
	}
	if r.MatchedCount == 0 {
		log.Printf("inserted one [ %v ] to collection [ %v ]", t, userChosenCollection.Name())
	} else if r.ModifiedCount == 1 {
		log.Printf("updated one [ %v ] to collection [ %v ]", t, userChosenCollection.Name())
	}

	return nil

}

func (db *MongoBongo) GetCollection(c string) *mongo.Collection {
	switch c {
	case db.CardsCollection:
		log.Printf("cards collection name is being set [ %v ]", c)
		return db.Client.Database(db.Database).Collection(db.CardsCollection)
	case db.SetsCollection:
		log.Printf("sets collection name is being set [ %v ]", c)
		return db.Client.Database(db.Database).Collection(db.SetsCollection)
	case db.UsersCollection:
		log.Printf("users collection name is being set [ %v ]", c)
		return db.Client.Database(db.Database).Collection(db.UsersCollection)
	default:
		log.Printf("default collection name is being set [ %v ]", c)
		return db.Client.Database(db.Database).Collection(db.CardsCollection)
	}
}

func InitDatabase(c config.Config) MongoBongo {

	client, err := mongo.NewClient(options.Client().ApplyURI(c.Mongo.Url))
	if err != nil {
		log.Fatalf("error init mongo client [%v]", err)
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("My name is carol, bith [%v]", err)
	}

	db := &MongoBongo{
		Client:          client,
		Database:        c.Mongo.Database,
		CardsCollection: c.Mongo.CardsCollection,
		SetsCollection:  c.Mongo.SetsCollection,
		UsersCollection: c.Mongo.UsersCollection,
	}

	return *db
}
func (db *MongoBongo) GetCardById(id string) (card *dto.Card) {

	c := db.Client.Database(db.Database).Collection(db.CardsCollection)

	log.Println("database  : ", db.Database)
	log.Println("collection: ", db.CardsCollection)
	log.Println("id        : ", id)

	resp := c.FindOne(context.Background(), bson.M{"id": id}).Decode(card)

	log.Printf("response from mongo [ %v ]", resp)
	return
}
