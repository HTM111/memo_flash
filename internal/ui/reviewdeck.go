package ui

import (
	"memoflash/internal/models"
	"memoflash/internal/utils"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/text/text"
	"cogentcore.org/core/tree"
)

type ReviewDeck struct {
	core.Frame
	deckData *models.Deck
}

func (deck *ReviewDeck) Init() {
	deck.Frame.Init()
	deck.Styler(func(s *styles.Style) {
		s.CenterAll()
		s.Background = colors.Scheme.SurfaceContainer
		s.Border.Radius.Set(units.Dp(12))
		s.Padding.SetAll(units.Dp(12))
		s.Gap.Set(units.Dp(20))
		s.Grow.Set(1, 0)
	})
	tree.AddChild(deck, func(w *core.Frame) {
		w.Styler(func(s *styles.Style) {
			s.Border.Radius = styles.BorderRadiusFull
			s.Min.Set(units.Dp(15))
		})
		w.Updater(func() {
			if deck.deckData != nil {
				w.Styler(func(s *styles.Style) {
					s.Background = colors.Uniform(CategoryColors[deck.deckData.CategoryIndex])
				})
			}
		})
	})

	tree.AddChild(deck, func(w *core.Frame) {
		w.Styler(func(s *styles.Style) {
			s.Direction = styles.Column
			s.Grow.Set(1, 0)
		})

		tree.AddChild(w, func(w *core.Text) {
			w.Styler(func(s *styles.Style) {
				s.Font.Size.Set(20, units.UnitDp)
			})
			w.Updater(func() {
				if deck.deckData != nil {
					w.SetText(deck.deckData.Title)
				}
			})
		})
		tree.AddChild(w, func(w *core.Text) {
			w.Styler(func(s *styles.Style) {
				s.Text.Align = text.Start
				s.Color = colors.Scheme.OnSurfaceVariant
			})
			w.Updater(func() {
				if deck.deckData != nil {
					w.SetText(deck.deckData.Description)
				}
			})
		})

		tree.AddChild(w, func(w *core.Frame) {
			w.Styler(func(s *styles.Style) {
				s.Direction = styles.Row
				s.Align.Items = styles.Center
				s.Gap.Set(units.Dp(5))
				s.Grow.Set(1, 0)
			})

			tree.AddChild(w, func(w *core.Icon) {
				w.SetIcon(icons.Schedule)
				w.Styler(func(s *styles.Style) {
					s.Color = colors.Scheme.OnSurfaceVariant
				})

			})

			tree.AddChild(w, func(w *core.Text) {
				w.Styler(func(s *styles.Style) {
					s.Color = colors.Scheme.OnSurfaceVariant
				})
				w.Updater(func() {
					if deck.deckData != nil {
						w.SetText(utils.FormatDuration(deck.deckData.LastStudied))
					}
				})
			})
		})
	})
}

func (deck *ReviewDeck) SetData(d *models.Deck) {
	deck.deckData = d
	deck.Update()
}
