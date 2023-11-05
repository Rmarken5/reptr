package logic

import (
	"context"
	"github.com/rmarken/reptr/internal/models"
	"time"
)

type (
	Controller interface {
		InsertDeck(ctx context.Context, deck models.Deck) error
		GetDecks(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]models.WithCards, error)
		AddCardToDeck(ctx context.Context, deckName string, card models.Card) error
		UpdateCard(ctx context.Context, card models.Card) error
	}
)
