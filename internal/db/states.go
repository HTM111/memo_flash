package db

import (
	"memoflash/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type CounterFilter struct {
	Condition any
	Table     string
}

func (database *Database) UpdateStats(set map[string]any) error {
	query := sq.Update("states")
	for key, value := range set {
		query = query.Set(key, value)
	}
	_, err := query.RunWith(database.db).Exec()
	return err
}
func (database *Database) SelectStats() (*models.States, error) {
	var lastUpdated int64
	queryBuilder := sq.Select("*").From("states")
	var stats = &models.States{}

	err := queryBuilder.RunWith(database.db).QueryRow().Scan(&stats.DayStreak, &lastUpdated)
	if lastUpdated != 0 {
		stats.LastTimeUpdated = time.Unix(lastUpdated, 0)
	}
	return stats, err
}
func (database *Database) Count(filter CounterFilter) (float32, error) {
	queryBuilder := sq.Select("COUNT(*)").From(filter.Table)
	if filter.Condition != nil {
		queryBuilder = queryBuilder.Where(filter.Condition)
	}
	var count float32
	err := queryBuilder.RunWith(database.db).QueryRow().Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
