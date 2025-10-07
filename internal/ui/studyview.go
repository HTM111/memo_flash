package ui

import (
	"fmt"
	"image/color"
	"strconv"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/text/rich"
	"cogentcore.org/core/tree"
)

func (app *App) createStudyTab(frame *core.Frame) {
	tree.AddChild(frame, func(container *core.Frame) {
		container.Styler(func(s *styles.Style) {
			s.Grow.Set(1, 1)
			s.Direction = styles.Column
			s.Margin.SetAll(units.Dp(15))
			s.Gap.Set(units.Dp(10))
		})

		app.makeStudyHeader(container)
		app.createStatsSection(container)
		app.createDecksSection(container)
	})
}

func (app *App) makeStudyHeader(parent *core.Frame) {
	tree.AddChild(parent, func(header *core.Frame) {
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

func (app *App) createStatsSection(parent *core.Frame) {
	tree.AddChildAt(parent, "stats-section", func(section *core.Frame) {
		section.Styler(func(s *styles.Style) {
			s.Grow.Set(1, 0)
			s.Direction = styles.Row
			s.Gap.Set(units.Dp(16))
		})
		app.makeStateCard(section, "due-cards", "Due", func() string {
			return fmt.Sprintf("%d", app.StudyManager.GetDueCount())
		}, icons.TodayFill, colors.Orange)
		app.makeStateCard(section, "day-streak", "Day Streak", func() string {
			return fmt.Sprintf("%d", app.StudyManager.GetDayStreak())
		}, icons.BoltFill, colors.Springgreen)
		app.makeStateCard(section, "progress", "Progress", func() string {
			return fmt.Sprintf("%d%%", app.StudyManager.GetProgress())
		}, icons.BarChart, colors.Mediumvioletred)
	})
}

func (app *App) makeStateCard(section *core.Frame, id, title string, value func() string, icon icons.Icon, color color.RGBA) {
	tree.AddChildAt(section, id, func(w *StatCard) {
		w.SetTitle(title)
		w.Updater(func() {
			if value != nil {
				w.SetValue(value())
			}
		})
		w.SetIcon(icon)
		w.SetColor(color)
	})
}

func (app *App) createDecksSection(parent *core.Frame) {
	tree.AddChildAt(parent, "decks-section", func(section *core.Frame) {
		section.Styler(func(s *styles.Style) {
			s.Grow.Set(1, 1)
			s.Direction = styles.Column
			s.Background = colors.Scheme.SurfaceContainerLow
			s.Border.Color.SetAll(colors.Scheme.SurfaceContainer)
			s.Border.Radius = styles.BorderRadiusMedium
			s.Padding.Set(units.Dp(16))
			s.Gap.Set(units.Dp(16))
		})

		tree.AddChild(section, func(header *core.Text) {
			header.SetText("Recently Studied").SetType(core.TextTitleLarge).Styler(func(s *styles.Style) {
				s.Font.Weight = rich.ExtraBold
			})
		})

		section.Maker(func(p *tree.Plan) {
			app.loadRecentDecks(p)
		})
	})
}

func (app *App) loadRecentDecks(p *tree.Plan) {
	recentlyStudied := app.StudyManager.GetRecentlyStudiedDecks()
	if len(recentlyStudied) == 0 {
		app.emptyState(p, "No decks have been studied yet", icons.School)
		return
	}

	for i, deck := range recentlyStudied {
		tree.AddAt(p, strconv.Itoa(i)+":"+strconv.Itoa(deck.ID), func(deckCard *ReviewDeck) {
			deckCard.SetData(deck)
		})
	}
}

func (app *App) emptyState(p *tree.Plan, message string, ic icons.Icon) {
	tree.Add(p, func(w *emptyState) {
		w.SetIcon(ic)
		w.SetMessage(message)
	})
}
