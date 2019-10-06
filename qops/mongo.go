package qops

import (
	"context"
	"fmt"
	"log"
	"pokemontcg-api-client/dto"
	"pokemontcg-api-client/web"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DataAccess interface {
	GetCollection() mongo.Collection
	Upsert(interface{}, string)
	Count()
	//Delete()
}

type MongoBongo struct {
	client     *mongo.Client
	collection string
	URL        string
	database   string
}

func (db *MongoBongo) Count(c string) {

	db.collection = c
	collection := db.GetCollection()
	cursor, err := collection.Find(context.Background(), bson.M{"set": "Unified Minds"})
	if err != nil {
		fmt.Println("cursor error")
	}
	var cards []dto.Card

	err = cursor.All(context.Background(), &cards)
	if err != nil {
		fmt.Println("decode error")
	}

	fmt.Printf("Set size: %d\n", len(cards))
}

func (db *MongoBongo) Upsert(t interface{}, c string) error {
	var filter bson.M
	var ok bool
	if _, ok = t.(dto.Card); ok {
		filter = bson.M{"id": t.(dto.Card).ID}
	} else {
		filter = bson.M{"code": t.(dto.Set).Code}
	}
	update := bson.M{"$set": t}
	db.collection = c
	collection := db.GetCollection()
	r, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("Failed to insert one to %s collection [ %v ]", collection.Name(), err)
	}
	if r.MatchedCount == 0 {
		if ok {
			log.Printf("Inserted One to %s Collection:: %v", collection.Name(), t.(dto.Card).ID)
		} else {
			log.Printf("Inserted One to %s Collection:: %v", collection.Name(), t.(dto.Set).Code)
		}

	} else if r.ModifiedCount == 1 {
		if ok {
			log.Printf("Updated One to %s Collection:: %v", collection.Name(), t.(dto.Card).ID)
		} else {
			log.Printf("Updated One to %s Collection:: %v", collection.Name(), t.(dto.Set).Code)
		}
	}

	return nil

}

func (db *MongoBongo) PopulateDB() {

	sets, err := web.GetAllSets()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range sets.Sets {
		log.Printf("Current Set: [ %s ]", v.Name)
		err := db.Upsert(v, "sets")
		if err != nil {
			log.Fatalf("Failed to insert set [ %v ]", err)
		}
		cards := web.GetCardsBySetCode(v.Code)

		if len(cards.Cards) != v.TotalCards {
			fmt.Printf("[WARNING] Did not receive all cards in set: %s - Actual: %d, Expected: %d\n", v.Name, len(cards.Cards), v.TotalCards)
		}

		for _, c := range cards.Cards {
			err := db.Upsert(c, "cards")
			if err != nil {
				log.Fatalf("Failed to insert cards [ %v ]", err)
			}
			// log.Printf("Inserted card [ result: %v ]", r)
		}

	}

}

func (db *MongoBongo) GetCollection() *mongo.Collection {
	return db.client.Database(db.database).Collection(db.collection)

}

func InitDatabase(url string, d string) *MongoBongo {

	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("My name is carol and i am bith")
	}
	ping := client.Ping(ctx, readpref.Primary())
	if ping != nil {
		fmt.Println("Couldn't ping DB ", ping)
	}

	db := &MongoBongo{
		client:   client,
		URL:      url,
		database: d,
	}

	return db
}

func setOps(url string) *options.ClientOptions {

	ops := options.Client().ApplyURI(url)
	// ops.SetSocketTimeout(30 * time.Second)
	// ops.SetAppName("Elite Trainer TCG Data Builder")
	// ops.SetConnectTimeout(5 * time.Second)
	return ops
}
