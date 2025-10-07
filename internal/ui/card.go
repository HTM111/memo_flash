package ui

import (
	"memoflash/internal/models"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/text/rich"
	"cogentcore.org/core/tree"
)

type Card struct {
	core.Frame
	Data     *models.Card
	onDelete func()
	onEdit   func()
}

func (card *Card) SetData(data *models.Card) {
	card.Data = data
}

func (card *Card) SetEdit(f func()) {
	card.onEdit = f
}
func (card *Card) SetDelete(f func()) {
	card.onDelete = f
}
func (card *Card) Init() {
	card.Frame.Init()
	card.Styler(func(s *styles.Style) {
		s.Grow.X = 1
		s.Direction = styles.Column
		s.Background = colors.Scheme.SurfaceContainer
		s.Border.Radius = styles.BorderRadiusSmall
		s.Padding.Set(units.Dp(17))

	})

	tree.AddChild(card, func(w *core.Text) {
		w.SetText("Front")
		w.Styler(func(s *styles.Style) {
			s.Font.Weight = rich.Bold
		})
	})
	tree.AddChild(card, func(w *core.Text) {
		w.Updater(func() {
			if card.Data != nil {
				w.SetText(card.Data.Front)
			}
		})
	})

	tree.AddChild(card, func(w *core.Text) {

		w.SetText("Back")
		w.Styler(func(s *styles.Style) {
			s.Font.Weight = rich.Bold
		})
	})
	tree.AddChild(card, func(w *core.Text) {
		w.Updater(func() {
			if card.Data != nil {
				w.SetText(card.Data.Back)
			}

		})
	})

	tree.AddChild(card, func(w *core.Frame) {
		w.Styler(func(s *styles.Style) {
			s.Margin.SetTop(units.Dp(10))
			s.Margin.SetRight(units.Dp(10))
			s.Grow.Set(1, 0)
			s.Align.Items = styles.Center
		})
		tree.AddChild(w, func(w *core.Button) {
			w.Styler(func(s *styles.Style) {
				s.Grow.Set(1, 0)

			})
			w.SetIcon(icons.Edit)
			w.SetText("Edit")
			w.OnClick(func(e events.Event) {
				if card.onEdit != nil {
					card.onEdit()
				}
			})
		})

		tree.AddChild(w, func(w *core.Button) {
			w.Styler(func(s *styles.Style) {
				s.Grow.Set(1, 0)

			})
			w.SetIcon(icons.Delete)
			w.SetText("Delete")
			w.OnClick(func(e events.Event) {
				if card.onDelete != nil {
					card.onDelete()
					card.onDelete = nil
				}
			})
		})
	})
}
