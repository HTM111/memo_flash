package ui

import (
	"memoflash/internal/models"

	"cogentcore.org/core/core"
)

type deckrepo interface {
	GetDecks() []*models.Deck
	AddDeck(*models.Deck)
	GetDeck(id int) *models.Deck
	DeleteDeck(id int)
}

func (app *App) AddDeck(d *models.Deck) {
	app.Decks = append(app.Decks, d)
	app.DeckMap[d.ID] = d

}
func (app *App) GetDecks() []*models.Deck {
	return app.Decks
}
func (app *App) DeleteDeck(id int) {
	for i, deckitem := range app.Decks {
		if deckitem.ID == id {
			app.Decks = append(app.Decks[:i], app.Decks[i+1:]...)
			delete(app.DeckMap, id)
			break
		}
	}
}
func (app *App) FetchDecks() {
	list, err := app.Services.GetDecks()
	if err != nil {
		core.ErrorDialog(app, err, "Error Getting Decks")
		return
	}
	if len(list) == 0 {
		list = []*models.Deck{}
	}
	for _, item := range list {
		app.DeckMap[item.ID] = item
	}
	app.Decks = list

}
func (app *App) GetDeck(id int) *models.Deck {
	if deck, found := app.DeckMap[id]; found {
		return deck
	}
	return nil

}
