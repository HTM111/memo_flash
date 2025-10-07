package db

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
)

type OrderBy string

const (
	DeckID                 OrderBy = "ID"
	DeckTitle              OrderBy = "Title"
	DeckDescription        OrderBy = "Description"
	DeckLastStudied        OrderBy = "LastStudied"
	DeckCategoryColorIndex OrderBy = "CategoryColorIndex"
	CreateAt               OrderBy = "CreatedAt"
	NONE                   OrderBy = "NONE"
)

type Database struct {
	db   *sql.DB
	psql sq.StatementBuilderType
}

func SetupDatabase(name string) (*Database, error) {
	db, err := sql.Open("sqlite3", name+"?_foreign_keys=on") // <-- add this
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return &Database{db: db, psql: sq.StatementBuilder.RunWith(db)}, nil
}

func (database *Database) InitSchema() error {

	tableQueries := []string{
		`CREATE TABLE IF NOT EXISTS states (
			dayStreak       INTEGER DEFAULT 0,
			lastTimeUpdated INTEGER DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS decks (
			ID                 INTEGER PRIMARY KEY AUTOINCREMENT,
			Title              TEXT,
			Description        TEXT,
			LastStudied        INTEGER,
			CategoryColorIndex TINYINT DEFAULT 0,
			CreatedAt          INTEGER DEFAULT (strftime('%s','now'))
		)`,
		`CREATE TABLE IF NOT EXISTS cards (
			ID           INTEGER PRIMARY KEY AUTOINCREMENT,
			Front        TEXT,
			Back         TEXT,
			Stability    REAL DEFAULT 1,
			Difficulty   REAL DEFAULT 0.3,
			LastStudied  INTEGER,
			ParentDeckId INTEGER,
			Interval     INTEGER,
			FOREIGN KEY (ParentDeckId) REFERENCES decks(ID) ON DELETE CASCADE
		)`,
	}
	for _, query := range tableQueries {
		if _, err := database.db.Exec(query); err != nil {
			return fmt.Errorf("create table: %w", err)
		}
	}
	sqlStmt, args, _ := sq.Select("dayStreak").From("states").Limit(1).ToSql()

	var dayStreak int
	err := database.db.QueryRow(sqlStmt, args...).Scan(&dayStreak)

	if err == sql.ErrNoRows {
		_, err = sq.Insert("states").
			Columns("dayStreak", "lastTimeUpdated").
			Values(0, 0).
			RunWith(database.db).
			Exec()
		if err != nil {
			return fmt.Errorf("seed initial state: %w", err)
		}
	}
	return nil
}
func (database *Database) Close() {
	database.db.Close()
}
