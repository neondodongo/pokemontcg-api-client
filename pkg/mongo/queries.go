package mongo

import (
	"context"
	"fmt"
	"net/url"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/dto"

	"github.com/rs/zerolog/log"
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
	Client   *mongo.Client
	Database string
}

// Upsert filters by interface type and attempts to upsert a corresponding document to MongoDB
func (db *MongoBongo) Upsert(t interface{}, filter bson.M, col string) error {
	update := bson.M{"$set": t}

	log.Printf("cards collection name is being set [ %v ]", col)

	c := db.Client.Database(db.Database).Collection(col)

	res, err := c.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("Failed to insert one to %s collection [ %v ]", c.Name(), err)
	}

	if res.MatchedCount == 0 {
		log.Info().Msgf("inserted one [ %v ] to collection [ %v ]", t, c.Name())
	} else if res.ModifiedCount == 1 {
		log.Info().Msgf("updated one [ %v ] to collection [ %v ]", t, c.Name())
	} else if res.ModifiedCount > 1 {
		log.Warn().Msgf("unexpected update [ %v ] to collection [ %v ]", t, c.Name())
	}

	return nil

}

// InitDatabase creates an instance of a MongoBongo
func InitDatabase(c config.Config) MongoBongo {

	client, err := mongo.NewClient(options.Client().ApplyURI(c.Mongo.Url))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create MongoDB client")
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize MongoDB client")
	}

	db := &MongoBongo{
		Client:   client,
		Database: c.Mongo.Database,
	}

	return *db
}

// GetCardById executes a query to retrieve a Card filtered by ID
func (db *MongoBongo) GetCardById(id string, col string) (card *dto.Card) {

	c := db.Client.Database(db.Database).Collection(col)

	log.Debug().Msg("querying with...")
	log.Debug().Msgf("database  : ", db.Database)
	log.Debug().Msgf("collection: ", col)
	log.Debug().Msgf("id        : ", id)

	resp := c.FindOne(context.Background(), bson.M{"id": id}).Decode(&card)

	log.Debug().Msgf("response from mongo [ %v ]", resp)

	return
}

func Filter(params url.Values) bson.M {
	var filter []bson.M

	if len(params) > 0 {
		for k, v := range params {
			filter = append(filter, bson.M{k: v[0]})
		}

		log.Debug().Msgf("Filter By: %v", filter)

		return bson.M{"$and": filter}
	}

	return bson.M{}
}

func (db *MongoBongo) Find(filter bson.M, col string) (interface{}, error) {
	// if reflect.TypeOf(&t).Kind() != reflect.Slice {
	// 	return errors.New("illegal query result target: target interface must be of type slice")
	// }

	c := db.Client.Database(db.Database).Collection(col)
	ctx := context.Background()

	cursor, err := c.Find(context.Background(), filter)
	if err != nil {
		log.Error().Err(err).Msgf("error finding documents using filter: %q", filter)
	}

	results := make([]interface{}, 0)
	if err = cursor.All(ctx, &results); err != nil {
		log.Error().Err(err).Msg("error decoding results")
	}

	return results, nil
}

// GetFilterCards creates a filter on Card attribues returning a slice of Card
func (db *MongoBongo) GetCards(filter bson.M, col string) []dto.Card {

	// options := options.Find().SetLimit(100)
	c := db.Client.Database(db.Database).Collection(col)

	cursor, err := c.Find(context.Background(), filter)
	if err != nil {
		log.Error().Err(err).Msgf("error finding documents using filter: %q", filter)
	}

	cards := []dto.Card{}

	for cursor.Next(context.Background()) {
		card := dto.Card{}

		if err := cursor.Decode(&card); err != nil {
			log.Error().Err(err).Msg("failed to decode Card")
		}

		cards = append(cards, card)
	}

	log.Debug().Int("results", len(cards))

	return cards
}

//GetFilterSets creates a filter on Set attributes returning a slice of Set
func (db *MongoBongo) GetFilterSets(filter bson.M, col string) []dto.Set {
	// options := options.Find().SetLimit(100)
	c := db.Client.Database(db.Database).Collection(col)

	cursor, err := c.Find(context.Background(), filter)
	if err != nil {
		log.Error().Err(err).Msgf("error finding documents using filter: %q", filter)
	}

	sets := []dto.Set{}

	for cursor.Next(context.Background()) {
		set := dto.Set{}

		if err := cursor.Decode(&sets); err != nil {
			log.Error().Err(err).Msg("failed to decode to Set")
		}

		sets = append(sets, set)
	}

	log.Debug().Int("results", len(sets))

	return sets
}

// FindUser find one User by a provided filter
func (db *MongoBongo) FindUser(filter bson.M) (u dto.User, err error) {
	col := "user"
	c := db.Client.Database(db.Database).Collection(col)

	r, err := c.Find(context.Background(), filter)
	if err != nil {
		return dto.User{}, err
	}

	if err := r.Decode(&u); err != nil {
		return dto.User{}, err
	}

	return
}
