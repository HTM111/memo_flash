package services

import (
	"memoflash/internal/db"
	"memoflash/internal/models"
	"time"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
)

type CardService interface {
	CreateCard(Front string, Back string, deckId int) error
	DeleteCard(id int) error
	CountDueCardsFromDeck(deckId int) (int, error)
	GetTotalCardsInDeck(deckId int) (int, error)
	GetDueCardsFromDeck(deckId int) ([]*models.Card, error)
	GetProgress() (int, error)
	GetStreak() (int, error)
	CountDueCards() (int, error)
	GetAllDueCards() ([]*models.Card, error)
	GetCardsFromDeck(deckId int) ([]*models.Card, error)
	EditCard(id int, Front string, Back string) error
}
type cardService struct {
	db *db.Database
}

func NewCardService(db *db.Database) *cardService {
	return &cardService{db: db}
}
func (cs *cardService) GetCards() ([]*models.Card, error) {
	return cs.db.GetCards(db.CardFilter{})
}
func (cs *cardService) GetAllDueCards() ([]*models.Card, error) {
	return cs.db.GetCards(db.CardFilter{
		Order: "interval ASC",
		Where: squirrel.Or{
			squirrel.Eq{"interval": nil},
			squirrel.LtOrEq{"date(interval,'unixepoch')": time.Now().Format("2006-01-02")}},
	})
}
func (ds *cardService) GetTotalCardsInDeck(deckid int) (int, error) {
	count, err := ds.db.Count(db.CounterFilter{
		Condition: sq.Eq{"ParentDeckId": deckid},
		Table:     "cards",
	})
	return int(count), err
}

func (cs *cardService) CreateCard(Front string, Back string, deckId int) error {
	return cs.db.CreateCard(Front, Back, deckId)
}

func (cs *cardService) DeleteCard(id int) error {
	return cs.db.DeleteCard(db.CardFilter{
		Where: sq.Eq{"ID": id},
	})
}
func (cs *cardService) GetCardsFromDeck(deckId int) ([]*models.Card, error) {
	return cs.db.GetCards(db.CardFilter{
		Where: sq.Eq{"ParentDeckId": deckId},
	})
}
func (cs *cardService) GetDueCardsFromDeck(deckId int) ([]*models.Card, error) {
	return cs.db.GetCards(db.CardFilter{
		Order: "interval ASC",
		Where: squirrel.And{
			squirrel.Eq{"ParentDeckId": deckId},
			squirrel.Or{
				squirrel.Eq{"interval": nil},
				squirrel.LtOrEq{"date(interval,'unixepoch')": time.Now().Format("2006-01-02")},
			},
		},
	})
}

func (cs *cardService) EditCard(id int, Front string, Back string) error {
	return cs.db.EditCard(Front, Back, id)
}

func (cs *cardService) isYesterdayStudied() (bool, error) {
	value, err := cs.db.Count(db.CounterFilter{
		Condition: sq.Eq{"date(LastStudied,'unixepoch')": time.Now().AddDate(0, 0, -1).Format("2006-01-02")},
		Table:     "decks",
	})
	return value == 1, err
}
func (cs *cardService) GetStreak() (int, error) {
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

func (cs *cardService) CountDueCards() (int, error) {
	counts, err := cs.db.Count(db.CounterFilter{
		Condition: squirrel.Or{
			squirrel.Eq{"interval": nil},
			squirrel.LtOrEq{"date(interval,'unixepoch')": time.Now().Format("2006-01-02")},
		},
		Table: "cards",
	})
	return int(counts), err
}
func (cs *cardService) GetProgress() (int, error) {
	totalCards, err := cs.db.Count(db.CounterFilter{
		Table: "cards",
	})
	if err != nil {
		return 0, err
	}
	progress, err := cs.db.Count(db.CounterFilter{
		Table:     "cards",
		Condition: squirrel.NotEq{"interval": nil},
	})
	if err != nil {
		return 0, err
	}
	var ProgressPercentage int
	if totalCards == 0 {
		ProgressPercentage = 0
	} else {
		ProgressPercentage = int((progress / totalCards) * 100)
	}
	return ProgressPercentage, nil
}

func (cs *cardService) CountDueCardsFromDeck(deckId int) (int, error) {
	total, err := cs.db.Count(db.CounterFilter{
		Condition: squirrel.And{
			squirrel.Eq{"ParentDeckId": deckId},
			squirrel.Or{
				squirrel.Eq{"interval": nil},
				squirrel.LtOrEq{"date(interval,'unixepoch')": time.Now().Format("2006-01-02")},
			},
		},
		Table: "cards",
	})
	return int(total), err
}
