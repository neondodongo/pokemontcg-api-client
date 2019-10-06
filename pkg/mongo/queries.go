package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"pokemontcg-api-client/dto"
	"pokemontcg-api-client/pkg/card"
	"pokemontcg-api-client/pkg/config"
)

type DataAccess interface {
	GetCollection() mongo.Collection
	Upsert(interface{}, string)
	//Delete()
}

type MongoBongo struct {
	Mgo     *mongo.Client
	Config config.Config
}

func (db *MongoBongo) Upsert(t interface{}, c string) error {
	var filter bson.M

	//determine interface type
	switch t.(type){
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

func (db *MongoBongo) GetCollection(c string) *mongo.Collection{

	switch c{
	case "cards":
		return db.Mgo.Database(db.Config.Mongo.Database).Collection(db.Config.Mongo.CardsCollection)
	case "sets":
		return db.Mgo.Database(db.Config.Mongo.Database).Collection(db.Config.Mongo.SetsCollection)
	case "users":
		return db.Mgo.Database(db.Config.Mongo.Database).Collection(db.Config.Mongo.UsersCollection)
	default:
		return db.Mgo.Database(db.Config.Mongo.Database).Collection(db.Config.Mongo.CardsCollection)
	}


}

func InitDatabase(c config.Config) *MongoBongo {

	client, err := mongo.NewClient(options.Client().ApplyURI(c.Mongo.Url))
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
		Mgo:   client,
		Config: c,
	}

	return db
}

func(db *MongoBongo) PopulateDatabase() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		sets, err := card.GetAllSets()
		if err != nil {
			log.Fatal(err)
		}

		for _, s := range sets.Sets {
			err := db.Upsert(s, "sets")
			if err != nil {
				log.Printf("Failed to insert set [ %v ] with error [ %v ]", s.Code, err)
			}
			cards := card.GetCardsBySetCode(db.Config, s.Code)

			if len(cards.Cards) != s.TotalCards {
				log.Printf("[WARNING] Did not receive all card in set: %s - Actual: %d, Expected: %d\n", s.Name, len(cards.Cards), s.TotalCards)
			}

			for _, c := range cards.Cards {
				err := db.Upsert(c, "card")
				if err != nil {
					log.Printf("Failed to insert card [ %v ] with error [ %v ]", c.ID, err)
				}
			}
		}
	})
}