package services

type Service struct {
	DeckService
	CardService
}

// type states struct {
// 	DeckService
// 	CardService
// }

// func (states *states) GetStreak() {
// 	stats, err := states.GetAllDueCards()
// 	if err != nil {
// 		return 0, err
// 	}
// 	var lastTimeUpdated = stats.LastTimeUpdated
// 	var streak = stats.DayStreak
// 	if lastTimeUpdated.Format("2006-01-02") == time.Now().Format("2006-01-02") {
// 		return stats.DayStreak, nil
// 	}
// 	isStudied, err := cs.isYesterdayStudied()
// 	if err != nil {
// 		return 0, err
// 	}

// 	if isStudied {
// 		streak = streak + 1
// 	} else {
// 		streak = 0
// 	}
// 	err = cs.db.UpdateStats(map[string]any{"dayStreak": streak, "lastTimeUpdated": time.Now().Unix()})
// 	return streak, err

// }
