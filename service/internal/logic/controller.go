package logic

import (
	"context"
	"github.com/rmarken/reptr/internal/database/card"
	"github.com/rmarken/reptr/internal/database/deck"
	"time"
)

type (
	Controller interface {
		InsertDeck(ctx context.Context, deck deck.Deck) error
		GetDecks(ctx context.Context, from time.Time, to *time.Time, limit, offset int) ([]deck.Deck, error)
		AddCardToDeck(ctx context.Context, deckName string, card card.Card) error
		UpdateCard(ctx context.Context, card card.Card) error
	}
)
