package ui

import (
	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/tree"
)

type State struct {
	core.Frame
	Title   string
	Message string
}

func (c *State) Init() {
	c.Frame.Init()
	c.Styler(func(s *styles.Style) {
		s.Gap.Set(units.Dp(20))
		s.Direction = styles.Column
		s.Border.Radius.SetAll(units.Dp(10))
		s.CenterAll()
		s.Padding.Set(units.Dp(50))
		s.Background = colors.Scheme.SurfaceContainerLow
	})
	if c.Title == "" {
		c.Title = "Session Complete!"
	}
	tree.AddChild(c, func(w *core.Text) {
		w.SetType(core.TextHeadlineSmall)
		w.Styler(func(s *styles.Style) {
			s.SetTextWrap(false)
		})
		w.Updater(func() {
			w.SetText(c.Title)
		})
	})

	tree.AddChild(c, func(w *core.Text) {
		w.SetType(core.TextHeadlineSmall)
		w.Styler(func(s *styles.Style) {
			s.SetTextWrap(false)
		})
		w.Updater(func() {
			w.SetText(c.Message)
		})
	})
}

func StatePage(parent core.Widget) *State {
	fr := core.NewFrame(parent)
	fr.Styler(func(s *styles.Style) {
		s.Grow.Set(1, 1)
		s.CenterAll()
	})
	card := tree.New[State](fr)
	card.Message = "Congratulations!"
	return card
}
