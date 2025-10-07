package ui

import (
	"image/color"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/text/rich"
	"cogentcore.org/core/tree"
)

type StatCard struct {
	core.Frame
	StateIcon icons.Icon
	Title     string
	Value     string
	Color     color.Color
}

func (sc *StatCard) Init() {
	sc.Frame.Init()

	sc.StateIcon = icons.AdFill
	sc.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.CenterAll()
		s.Background = colors.Scheme.SurfaceContainerLow
		s.Border.Radius.Set(units.Dp(12))
		s.Grow.Set(1, 1)
		s.Padding.SetAll(units.Dp(16))
	})

	tree.AddChild(sc, func(iconContainer *core.Frame) {
		iconContainer.Styler(func(s *styles.Style) {
			s.Background = colors.Uniform(colors.ApplyOpacity(sc.Color, 0.10))
			s.Border.Radius = styles.BorderRadiusFull
			s.Padding.Set(units.Dp(16))
			s.CenterAll()
		})
		tree.AddChild(iconContainer, func(icon *core.Icon) {

			icon.SetIcon(sc.StateIcon)
			icon.Styler(func(s *styles.Style) {

				s.Color = colors.Uniform(sc.Color)
				s.Font.Size.Dp(25)
			})
		})

	})
	tree.AddChild(sc, func(txt *core.Text) {
		txt.Updater(func() {
			txt.SetText(sc.Value)

		})
		txt.Styler(func(s *styles.Style) {
			s.Color = colors.Uniform(colors.White)
			s.Font.Weight = rich.Bold
			s.Font.Size.Dp(24)
		})
	})

	tree.AddChild(sc, func(lbl *core.Text) {
		lbl.SetText(sc.Title)
		lbl.SetType(core.TextBodyMedium)
		lbl.Styler(func(s *styles.Style) {
			s.Color = colors.Uniform(colors.White)
		})

	})
}
func (sc *StatCard) SetIcon(icon icons.Icon) {
	sc.StateIcon = icon
}
func (sc *StatCard) SetTitle(title string) {
	sc.Title = title
}

func (sc *StatCard) SetValue(value string) {
	sc.Value = value
}

func (sc *StatCard) SetColor(color color.Color) {
	sc.Color = color
}
