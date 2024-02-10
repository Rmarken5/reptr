package pages

import "github.com/a-h/templ"

type (
	DeckViewPageData struct {
		DeckName string
		DeckID   string
		Content  templ.Component
	}
)
