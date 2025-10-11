package db

import (
	"database/sql"
	"fmt"
	"log"
	"memoflash/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type DeckFilter struct {
	OrderBy string
	Limit   uint64
	Where   any
}

func (database *Database) GetDecks(filter DeckFilter) ([]*models.Deck, error) {
	var decks []*models.Deck

	SelectBuilder := sq.Select(
		"decks.ID",
		"decks.Title",
		"decks.Description",
		"decks.LastStudied",
		"decks.CategoryColorIndex",
		"decks.CreatedAt",
		"COALESCE(COUNT(cards.ID), 0) as total_cards",
		"COALESCE((SELECT COUNT(*) FROM cards WHERE cards.ParentDeckId = decks.ID AND (cards.interval IS NULL OR date(cards.interval,'unixepoch') <= date('now'))), 0) as due_cards").
		From("decks").
		LeftJoin("cards ON decks.ID = cards.ParentDeckId").
		GroupBy("decks.ID", "decks.Title", "decks.Description", "decks.LastStudied", "decks.CategoryColorIndex", "decks.CreatedAt")

	if filter.Limit > 0 {
		SelectBuilder = SelectBuilder.Limit(filter.Limit)
	}
	if filter.OrderBy != "" {
		SelectBuilder = SelectBuilder.OrderBy(filter.OrderBy)
	}
	if filter.Where != nil {
		SelectBuilder = SelectBuilder.Where(filter.Where)
	}

	rows, err := SelectBuilder.RunWith(database.db).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		deck := new(models.Deck)
		var LastStudiedinInt sql.NullInt64
		var CreatedAt sql.NullInt64
		var TotalCards int
		var DueCards int

		err = rows.Scan(&deck.ID, &deck.Title, &deck.Description, &LastStudiedinInt, &deck.CategoryIndex, &CreatedAt, &TotalCards, &DueCards)

		if err != nil {
			log.Println("Error Scanning Row", err)
			continue
		}

		if CreatedAt.Valid {
			deck.CreatedAt = time.Unix(CreatedAt.Int64, 0)
		}
		if LastStudiedinInt.Valid {
			deck.LastStudied = time.Unix(LastStudiedinInt.Int64, 0)
		}

		deck.TotalCards = TotalCards
		deck.DueCards = DueCards

		decks = append(decks, deck)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return decks, nil
}
func (database *Database) CreateDeck(Title string, Description string, CategoryColorIndex int) (int, error) {
	currentTime := time.Now().Unix()
	result, err := sq.Insert("decks").Columns(
		"Title", "Description", "CategoryColorIndex", "CreatedAt",
	).Values(Title, Description, CategoryColorIndex, currentTime).RunWith(database.db).Exec()
	if err != nil {
		return 0, err
	}
	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastInsertedId), nil
}
func (database *Database) EditDeck(id int, title, description string, categoryIndex int) error {
	result, err := sq.Update("decks").
		Set("Title", title).Set("Description", description).Set("CategoryColorIndex", categoryIndex).
		Where(sq.Eq{"id": id}).RunWith(database.db).Exec()
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		log.Println("No rows affected", id)
	}

	return nil
}

func (database *Database) DeleteDeck(filter DeckFilter) error {
	query := sq.Delete("decks").RunWith(database.db)
	if filter.Where != nil {
		query = query.Where(filter.Where)
	}
	result, err := query.Exec()
	if err != nil {
		return fmt.Errorf("Error executing query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		log.Println("No rows affected")
	}
	return nil
}
func (database *Database) UpdateReadTime(id int) error {
	timeInterval := time.Now().Unix()
	_, err := sq.Update("decks").Set("LastStudied", timeInterval).Where(sq.Eq{"id": id}).RunWith(database.db).Exec()
	if err != nil {
		return fmt.Errorf("Error preparing statement: %w", err)
	}
	return nil
}
