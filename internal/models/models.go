package models

import (
	"image/color"
	"time"

	"cogentcore.org/core/icons"
)

// `CREATE TABLE IF NOT EXISTS states (
// 			dayStreak       INTEGER DEFAULT 0,
// 			lastTimeUpdated INTEGER DEFAULT 0
// 		)`,

type States struct {
	DayStreak       int       `db:"dayStreak"`
	LastTimeUpdated time.Time `db:"lastTimeUpdated"`
}

type Deck struct {
	ID            int
	Title         string
	Description   string
	CategoryIndex int
	LastStudied   time.Time
	finishedCards int
	TotalCards    int
	DueCards      int
	CreatedAt     time.Time
}

type StatsCard struct {
	StateIcon icons.Icon
	Title     string
	Value     func() string
	Color     color.Color
}
type Card struct {
	ID           int       `db:"ID"`
	Front        string    `db:"Front"`
	Back         string    `db:"Back"`
	ParentDeckId int       `db:"ParentDeckId"`
	LastStudied  time.Time `db:"LastStudied"`
	Stability    float64   `db:"Stability"`
	Difficulty   float64   `db:"Difficulty"`
	Interval     time.Time `db:"Interval"`
}

func (card *Card) IsDue() bool {
	return card.Interval.IsZero() || card.Interval.Before(time.Now())
}
