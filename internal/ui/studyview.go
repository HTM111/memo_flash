package ui

import (
	"fmt"
	"memoflash/internal/models"
	"memoflash/internal/services"
	"sort"
	"strconv"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/text/rich"
	"cogentcore.org/core/tree"
)

type StudyTab struct {
	core.Frame
	deckrepo  deckrepo
	services  *services.Service
	Due       int
	DayStreak int
	Progress  int
}

func (st *StudyTab) Init() {
	st.Frame.Init()
	st.Styler(func(s *styles.Style) {
		s.Grow.Set(1, 1)
		s.Direction = styles.Column
		s.Margin.SetAll(units.Dp(15))
		s.Gap.Set(units.Dp(10))
	})

	st.makeStudyHeader()
	st.createStatsSection()
	st.createDecksSection()
}
func (st *StudyTab) fetchStats() error {
	due, err := st.services.CountDueCards()
	if err != nil {
		return err
	}
	st.Due = due
	daystreak, err := st.services.GetStreak()
	if err != nil {
		return err
	}
	st.DayStreak = daystreak
	progress, err := st.services.GetProgress()
	if err != nil {
		return err
	}
	st.Progress = progress
	return nil

}
func (st *StudyTab) makeStudyHeader() {
	tree.AddChild(st, func(header *core.Frame) {
		header.Styler(func(s *styles.Style) {
			s.Direction = styles.Column
			s.Gap.Set(units.Dp(8))
		})

		tree.AddChild(header, func(title *core.Text) {
			title.SetText("Study").SetType(core.TextHeadlineLarge).Styler(func(s *styles.Style) {
				s.Font.Weight = rich.Bold
			})
		})

		tree.AddChild(header, func(subtitle *core.Text) {
			subtitle.SetText("Track your progress and continue learning")
			subtitle.SetType(core.TextBodyLarge)

			subtitle.Styler(func(s *styles.Style) {
				s.Color = colors.Scheme.OnSurfaceVariant
				s.SetTextWrap(false)
			})
		})
	})
}

func (st *StudyTab) createStatsSection() {

	tree.AddChild(st, func(section *core.Frame) {
		section.Styler(func(s *styles.Style) {
			s.Grow.Set(1, 0)
			s.Direction = styles.Row
			s.Gap.Set(units.Dp(16))
		})

		tree.AddChildAt(section, "due-today-card", func(w *StatCard) {
			w.SetTitle("Due")
			w.Updater(func() {
				w.SetValue(fmt.Sprintf("%d", st.Due))

			})
			w.SetIcon(icons.TodayFill)
			w.SetColor(colors.Orange)
		})
		tree.AddChildAt(section, "day-streak-card", func(w *StatCard) {
			w.SetTitle("Day Streak")
			w.Updater(func() {
				w.SetValue(fmt.Sprintf("%d", st.DayStreak))
			})
			w.SetIcon(icons.BoltFill)
			w.SetColor(colors.Springgreen)
		})
		tree.AddChildAt(section, "progress", func(w *StatCard) {
			w.SetTitle("Progress")
			w.Updater(func() {
				w.SetValue(fmt.Sprintf("%d%%", st.Progress))

			})
			w.SetIcon(icons.BarChart)
			w.SetIcon(icons.BoltFill)
			w.SetColor(colors.Mediumvioletred)
		})
		section.OnFirst(events.Show, func(e events.Event) {
			err := st.fetchStats()
			if err != nil {
				core.ErrorSnackbar(st, err, "Error fetching Stats")
				return
			}
			section.Update()
		})
	})

}
func (st *StudyTab) GetRecentDecks() []*models.Deck {
	sortedList := make([]*models.Deck, 0)
	for _, item := range st.deckrepo.GetDecks() {
		if !item.LastStudied.IsZero() {
			sortedList = append(sortedList, item)
		}
	}
	sort.Slice(sortedList, func(i, j int) bool {
		return sortedList[i].LastStudied.After(sortedList[j].LastStudied)
	})
	if len(sortedList) > 3 {
		sortedList = sortedList[:3]
		return sortedList
	}
	return sortedList

}
func (st *StudyTab) createDecksSection() {
	tree.AddChild(st, func(header *core.Text) {
		header.SetText("Recently Studied").SetType(core.TextTitleLarge).Styler(func(s *styles.Style) {
			s.Font.Weight = rich.ExtraBold
		})
	})

	tree.AddChildAt(st, "decks-section", func(section *core.Frame) {
		section.Styler(func(s *styles.Style) {
			s.Grow.Set(1, 1)
			s.Direction = styles.Column
			s.Background = colors.Scheme.SurfaceContainerLow
			s.Border.Color.SetAll(colors.Scheme.SurfaceContainer)
			s.Border.Radius = styles.BorderRadiusMedium
			s.Padding.Set(units.Dp(16))
			s.Gap.Set(units.Dp(16))
		})
		section.Maker(func(p *tree.Plan) {
			if len(st.deckrepo.GetDecks()) == 0 {
				EmptyState(p, "No decks have been studied yet", icons.School)
			} else {
				for _, deck := range st.GetRecentDecks() {
					deckID := deck.ID
					tree.AddAt(p, strconv.Itoa(deckID), func(deckCard *RecentDeck) {
						deckCard.Updater(func() {
							deckCard.SetData(deck)
						})
					})
				}
			}
		})
	})
}
