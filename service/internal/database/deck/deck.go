package deck

import (
	"context"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	DataAccess interface {
		InsertDeck(ctx context.Context, deck Deck) error
		InsertCard(ctx context.Context, deckName string, card Card) error
	}

	Deck struct {
		Name  string `bson:"name"`
		Cards []Card `bson:"cards"`
	}
	CardType int
	Card     struct {
		Front string   `bson:"front"`
		Back  string   `bson:"back"`
		Kind  CardType `bson:"kind"`
	}
	DAO struct {
		collection *mongo.Collection
		log        zerolog.Logger
	}
)

const (
	BasicCard = iota
	MultipleChoice
)

func NewDataAccess(db *mongo.Database, log zerolog.Logger) *DAO {
	logger := log.With().Str("module", "DAO").Logger()
	collection := db.Collection("decks")
	return &DAO{
		collection: collection,
		log:        logger,
	}
}

func (c CardType) String() string {
	switch c {
	case BasicCard:
		return "basic"
	case MultipleChoice:
		return "multiple choice"
	}
	return "unknown"
}

func (d *DAO) InsertDeck(ctx context.Context, deck Deck) error {
	res, err := d.collection.InsertOne(ctx, deck)
	if err != nil {
		return err
	}

	d.log.Printf("%+v", res)
	return nil
}

func (d *DAO) InsertCard(ctx context.Context, deckName string, card Card) error {
	//TODO implement me
	panic("implement me")
}
