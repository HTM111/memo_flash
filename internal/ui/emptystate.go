package ui

import (
	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/tree"
)

type emptyState struct {
	core.Frame
	Icon        icons.Icon
	IconSize    float32
	MessageSize float32
	Message     string
}

func (nf *emptyState) Init() {
	nf.Frame.Init()

	nf.Icon = icons.Info
	nf.IconSize = 60
	nf.MessageSize = 20
	nf.Message = "Not found"

	nf.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.CenterAll()
		s.Min.Y.Set(300, units.UnitDp)
		s.Grow.Set(1, 1)
	})

	tree.AddChild(nf, func(w *core.Icon) {
		w.Styler(func(s *styles.Style) {
			s.Color = colors.Scheme.OnSurfaceVariant
			s.Font.Size.Set(nf.IconSize, units.UnitDp)
		})
		w.Updater(func() {
			w.SetIcon(nf.Icon)
		})
	})

	tree.AddChild(nf, func(w *core.Text) {
		w.Styler(func(s *styles.Style) {
			s.Color = colors.Scheme.OnSurfaceVariant
			s.Font.Size.Set(nf.MessageSize, units.UnitDp)
		})
		w.Updater(func() {
			w.SetText(nf.Message)
		})
	})
}

func (nf *emptyState) SetIconSize(size float32) {
	nf.IconSize = size
}

func (nf *emptyState) SetMessageSize(size float32) {
	nf.MessageSize = size
}

func (nf *emptyState) SetIcon(icon icons.Icon) {
	nf.Icon = icon
}

func (nf *emptyState) SetMessage(message string) {
	nf.Message = message
}
func EmptyState(p *tree.Plan, message string, ic icons.Icon) {
	tree.Add(p, func(w *emptyState) {
		w.SetIcon(ic)
		w.SetMessage(message)
	})
}
