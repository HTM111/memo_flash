package ui

import (
	"fmt"
	"log"
	"memoflash/internal/models"
	"memoflash/internal/services"
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
	services  *services.Service
	decks     []*models.Deck
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
			w.SetIcon(icons.BoltFill)
			w.SetColor(colors.Springgreen)
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
		section.OnFirst(events.Show, func(e events.Event) {
			recentlyStudied, err := st.services.GetRecentlyStudiedDecks()
			if err != nil {
				core.ErrorSnackbar(st, err, "Error retrieving recently studied decks")
				return
			}
			st.decks = recentlyStudied
			for _, d := range recentlyStudied {
				if d.LastStudied.IsZero() {
					log.Println("time is null", d.Title)
				}
			}
			section.Update()

		})
		section.Maker(func(p *tree.Plan) {
			if len(st.decks) == 0 {
				EmptyState(p, "No decks have been studied yet", icons.School)
				return
			}
			for i, deck := range st.decks {
				tree.AddAt(p, strconv.Itoa(i)+":"+strconv.Itoa(deck.ID), func(deckCard *ReviewDeck) {
					deckCard.SetData(deck)
				})
			}
		})
	})
}

func EmptyState(p *tree.Plan, message string, ic icons.Icon) {
	tree.Add(p, func(w *emptyState) {
		w.SetIcon(ic)
		w.SetMessage(message)
	})
}
