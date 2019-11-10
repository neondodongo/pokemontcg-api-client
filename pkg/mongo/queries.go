package mongo

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/dto"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// Upsert filters by interface type and attempts to upsert a corresponding document to MongoDB
func (db *MongoBongo) Upsert(t interface{}) error {
	var filter bson.M
	var c string
	//determine interface type
	switch t.(type) {
	case dto.Card:
		filter = bson.M{"id": t.(dto.Card).ID}
		c = db.CardsCollection
	case dto.Set:
		filter = bson.M{"code": t.(dto.Set).Code}
		c = db.SetsCollection
	case dto.User:
		filter = bson.M{"username": t.(dto.User).Username}
		c = db.UsersCollection
	}

	update := bson.M{"$set": t}
	collection := db.SetCollection(c)
	r, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("Failed to insert one to %s collection [ %v ]", collection.Name(), err)
	}
	if r.MatchedCount == 0 {
		log.Printf("inserted one [ %v ] to collection [ %v ]", t, collection.Name())
	} else if r.ModifiedCount == 1 {
		log.Printf("updated one [ %v ] to collection [ %v ]", t, collection.Name())
	} // TODO: if modified count > 1, attempt delete duplicate records

	return nil

}

// SetCollection sets an active collection to a MongoBongo
func (db *MongoBongo) SetCollection(c string) *mongo.Collection {
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

// InitDatabase creates an instance of a MongoBongo
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

// GetCardById executes a query to retrieve a Card filtered by ID
func (db *MongoBongo) GetCardById(id string) (card *dto.Card) {

	c := db.Client.Database(db.Database).Collection(db.CardsCollection)

	log.Println("database  : ", db.Database)
	log.Println("collection: ", db.CardsCollection)
	log.Println("id        : ", id)

	resp := c.FindOne(context.Background(), bson.M{"id": id}).Decode(&card)

	log.Printf("response from mongo [ %v ]", resp)
	return
}

// GetFilterCards creates a filter on Card attribues returning a slice of Card
func (db *MongoBongo) GetFilterCards(params url.Values) []dto.Card {

	var filters []bson.M
	var filter bson.M

	if len(params) > 0 {
		for k, v := range params {
			filters = append(filters, bson.M{k: v[0]})

		}
		fmt.Println("filters: ", filters)
		filter = bson.M{"$and": filters}
	} else {
		filter = bson.M{}
	}

	// options := options.Find().SetLimit(100)
	c := db.Client.Database(db.Database).Collection(db.CardsCollection)

	cursor, err := c.Find(context.Background(), filter)
	if err != nil {
		log.Printf("error finding documents: %v", err)
	}

	cards := []dto.Card{}

	for cursor.Next(context.Background()) {
		card := dto.Card{}
		if err := cursor.Decode(&card); err != nil {
			log.Printf("Unable to decode card [%v]", err)
		}
		cards = append(cards, card)
	}

	log.Printf("Total cards from Filtered Search: %d", len(cards))

	return cards

}

//GetFilterSets creates a filter on Set attributes returning a slice of Set
func (db *MongoBongo) GetFilterSets(params url.Values) []dto.Set {

	var filters []bson.M
	var filter bson.M
	if len(params) > 0 {
		for k, v := range params {
			b, e := strconv.ParseBool(v[0])
			if e == nil {
				filters = append(filters, bson.M{k: b})
				continue
			}
			i, e := strconv.Atoi(v[0])
			if e == nil {
				filters = append(filters, bson.M{k: i})
				continue
			}
			filters = append(filters, bson.M{k: v[0]})
		}
		fmt.Println("filters: ", filters)
		filter = bson.M{"$and": filters}
	} else {
		filter = bson.M{}
	}

	// options := options.Find().SetLimit(100)
	c := db.Client.Database(db.Database).Collection(db.SetsCollection)
	cursor, err := c.Find(context.Background(), filter)
	if err != nil {
		log.Printf("error finding documents: %v", err)
	}

	sets := []dto.Set{}

	for cursor.Next(context.Background()) {
		set := dto.Set{}
		if err := cursor.Decode(&set); err != nil {
			log.Printf("Unable to decode set [%v]", err)
		}
		sets = append(sets, set)
	}

	log.Printf("Total sets from Filtered Search: %d", len(sets))

	return sets

}

// FindUserByUsername find one User by their username
func (db *MongoBongo) FindUserByUsername(un string) (u dto.User, err error) {
	filter := bson.M{"username": un}
	col := db.SetCollection(db.UsersCollection)
	r := col.FindOne(context.Background(), filter)
	if err = r.Decode(&u); err != nil {
		return
	}
	return
}

// FindUserByEmail find one User by their email address
func (db *MongoBongo) FindUserByEmail(em string) (u dto.User, err error) {
	filter := bson.M{"email": em}
	col := db.SetCollection(db.UsersCollection)
	r := col.FindOne(context.Background(), filter)
	if err = r.Decode(&u); err != nil {
		return
	}
	return
}
