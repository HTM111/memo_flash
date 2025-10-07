package db

import (
	"database/sql"
	"fmt"
	"log"
	"memoflash/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type CardFilter struct {
	Where any
	Order string
}

func (database *Database) GetCards(filter CardFilter) ([]*models.Card, error) {
	var cards []*models.Card
	queryBuilder := sq.Select("*").From("cards").RunWith(database.db)
	if filter.Where != nil {
		queryBuilder = queryBuilder.Where(filter.Where)
	}
	if filter.Order != "" {
		queryBuilder = queryBuilder.OrderBy(filter.Order)
	}
	rows, err := queryBuilder.Query()
	if err != nil {
		return cards, err
	}
	defer rows.Close()
	for rows.Next() {
		var interval sql.NullInt64
		var lastStudied sql.NullInt64
		card := new(models.Card)
		err = rows.Scan(&card.ID, &card.Front, &card.Back, &card.Stability, &card.Difficulty, &lastStudied, &card.ParentDeckId, &interval)
		if err != nil {
			log.Println("Error Scanning Row :", err)
			continue
		}
		if lastStudied.Valid && lastStudied.Int64 != 0 {
			validLastStudied := time.Unix(lastStudied.Int64, 0)
			card.LastStudied = validLastStudied
		}
		if interval.Valid && interval.Int64 != 0 {
			validInterval := time.Unix(interval.Int64, 0)
			card.Interval = validInterval
		}
		cards = append(cards, card)

	}
	return cards, nil
}
func (database *Database) CreateCard(front, back string, parentDeckId int) error {
	card, err := sq.Insert("cards").Columns("Front", "Back", "ParentDeckId").
		Values(front, back, parentDeckId).
		RunWith(database.db).
		Exec()
	if err != nil {
		return fmt.Errorf("Error Executing Statement: %w", err)
	}
	rowsAffected, err := card.RowsAffected()
	if err != nil {
		log.Println("Error Retrieving Rows Affected:", err)
	}
	if rowsAffected == 0 {
		log.Println("No rows affected")
	}
	return nil
}

func (database *Database) EditCard(front, back string, id int) error {

	_, err := sq.Update("cards").Set("Front", front).Set("Back", back).Where("id = ?", id).RunWith(database.db).Exec()

	if err != nil {
		return err
	}
	return nil
}

func (database *Database) DeleteCard(cardfilter CardFilter) error {
	query := sq.Delete("cards")
	if cardfilter.Where != nil {
		query = query.Where(cardfilter.Where)
	}
	result, err := query.RunWith(database.db).Exec()
	if err != nil {
		return fmt.Errorf("Error Executing Statement: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error Retrieving Rows Affected:", err)
	}
	if rowsAffected == 0 {
		log.Println("No rows affected")
	}

	return nil
}

func (database *Database) UpdateInterval(interval time.Time, Stability float64, Difficulty float64, id int) error {
	timeInterval := interval.Unix()
	s, err := database.db.Prepare("UPDATE cards SET interval = ?, stability = ?, difficulty = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = s.Exec(timeInterval, Stability, Difficulty, id)
	if err != nil {
		return fmt.Errorf("Error Executing Query: %w", err)
	}
	return nil
}
