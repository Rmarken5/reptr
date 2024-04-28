package pages

import "github.com/a-h/templ"

type (
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
