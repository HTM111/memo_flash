package services

import (
	"memoflash/internal/db"
	"memoflash/internal/models"

	sq "github.com/Masterminds/squirrel"

	"time"
)

type DeckService interface {
	DeleteDeck(id int) error
	CreateDeck(name string, description string, CategoryColorIndex int) (int, error)
	EditDeck(id int, name string, description string, CategoryColorIndex int) error
	GetDecks() ([]*models.Deck, error)
	GetRecentlyStudiedDecks() ([]*models.Deck, error)
	GetCardsFromDeck(deckId int) ([]*models.Card, error)
	UpdateInterval(id int, interval time.Time, Difficulty float64, Stability float64) error
	UpdateReadTime(id int) error
}

type deckService struct {
	db *db.Database
}

func NewDeckService(db *db.Database) DeckService {
	return &deckService{db: db}
}

func (ds *deckService) DeleteDeck(id int) error {
	return ds.db.DeleteDeck(db.DeckFilter{
		Where: sq.Eq{"ID": id},
	})
}
func (ds *deckService) CreateDeck(name string, description string, CategoryColorIndex int) (int, error) {
	return ds.db.CreateDeck(name, description, CategoryColorIndex)
}

func (ds *deckService) GetRecentlyStudiedDecks() ([]*models.Deck, error) {
	return ds.db.GetDecks(db.DeckFilter{
		Limit: 3,
		Where: sq.NotEq{
			"Interval": nil,
		},
		OrderBy: "decks.LastStudied DESC",
	})
}
func (ds *deckService) GetDecks() ([]*models.Deck, error) {
	d, err := ds.db.GetDecks(db.DeckFilter{
		OrderBy: "decks.LastStudied ASC",
	})
	return d, err
}
func (ds *deckService) UpdateReadTime(id int) error {
	return ds.db.UpdateReadTime(id)
}

func (ds *deckService) EditDeck(id int, name string, description string, CategoryColorIndex int) error {
	return ds.db.EditDeck(id, name, description, CategoryColorIndex)
}
func (ds *deckService) GetCardsFromDeck(deckId int) ([]*models.Card, error) {
	return ds.db.GetCards(db.CardFilter{
		Where: sq.Eq{"ParentDeckId": deckId},
	})

}
func (ds *deckService) UpdateInterval(id int, interval time.Time, Difficulty float64, Stability float64) error {
	return ds.db.UpdateInterval(interval, Difficulty, Stability, id)
}
func (cs *deckService) isYesterdayStudied() (bool, error) {
	value, err := cs.db.Count(db.CounterFilter{
		Condition: sq.Eq{"date(LastStudied,'unixepoch')": time.Now().AddDate(0, 0, -1).Format("2006-01-02")},
		Table:     "decks",
	})
	return value == 1, err
}
func (cs *deckService) GetStreak() (int, error) {
	stats, err := cs.db.SelectStats()
	if err != nil {
		return 0, err
	}
	var lastTimeUpdated = stats.LastTimeUpdated
	var streak = stats.DayStreak
	if lastTimeUpdated.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		return stats.DayStreak, nil
	}
	isStudied, err := cs.isYesterdayStudied()
	if err != nil {
		return 0, err
	}

	if isStudied {
		streak = streak + 1
	} else {
		streak = 0
	}
	err = cs.db.UpdateStats(map[string]any{"dayStreak": streak, "lastTimeUpdated": time.Now().Unix()})
	return streak, err
}
