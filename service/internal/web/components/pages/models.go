package pages

import (
	"github.com/a-h/templ"
	"github.com/rmarken/reptr/service/internal/web/components/dumb"
)

type (
	HomeData struct {
		Username string
		Groups   []HomeGroupData
		Decks    []dumb.Deck
	}

	HomeGroupData struct {
		ID        string
		GroupName string
		NumDecks  int
		NumUsers  int
	}

	DeckViewPageData struct {
		DeckName string
		DeckID   string
		Content  templ.Component
	}
	ErrorPageData struct {
		StatusCode string
		Status     string
		Error      string
		Msg        string
	}
)
