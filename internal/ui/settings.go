package ui

import (
	"slices"

	"cogentcore.org/core/base/iox/tomlx"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/tree"
)

type ParameterOption struct {
	core.Frame
	Title string
}

func (p *ParameterOption) Init() {
	p.Frame.Init()
	p.Styler(func(s *styles.Style) {
		s.Grow.Set(1, 0)
		s.Align.Items = styles.Center
		s.Gap.Set(units.Dp(8))
	})

}
func (p *ParameterOption) makeChooser(title string, options []string, defaultValue *string, onSelect func(s string)) {

	tree.AddChild(p, func(c *core.Text) {
		c.Styler(func(st *styles.Style) {
			st.SetTextWrap(false)
		})
		c.SetText(title)
	})
	tree.AddChild(p, func(w *core.Stretch) {

	})
	tree.AddChild(p, func(c *core.Chooser) {
		c.Styler(func(s *styles.Style) {
			s.Font.Size.Dp(15)
		})
		c.Updater(func() {
			c.SetCurrentIndex(slices.Index(options, *defaultValue))

		})
		c.SetStrings(options...)
		c.OnChange(func(e events.Event) {
			if onSelect != nil {
				onSelect(c.CurrentItem.GetText())
			}
		})
	})

}
func (p *ParameterOption) makeSpinner(title string, max, min, step float32) {

	tree.AddChild(p, func(c *core.Text) {
		c.Styler(func(st *styles.Style) {
			st.SetTextWrap(false)
		})
		c.SetText(title)
	})
	tree.AddChild(p, func(w *core.Stretch) {

	})
	tree.AddChild(p, func(c *core.Spinner) {
		c.Styler(func(s *styles.Style) {
			s.Font.Size.Dp(15)
		})
		c.SetMax(max)
		c.SetMin(min)
		c.SetStep(step)

	})

}

type AppSettings struct {
	core.SettingsBase

	DailyCardLimit int

	ThemeMode string
	CardSize  string
}

func (s *AppSettings) Defaults() {
	s.DailyCardLimit = 50
	s.ThemeMode = "Dark"
	s.CardSize = "Large"
}

func (s *AppSettings) Apply() {
	core.AppearanceSettings.Theme = getThemeFromText(s.ThemeMode)
	core.AppearanceSettings.Apply()
}
func (s *AppSettings) Save() error {
	return tomlx.Save(s, s.Filename())
}

func (s *AppSettings) Open() error {
	return tomlx.Open(s, s.Filename())
}
func getThemeFromText(s string) core.Themes {
	switch s {
	case "Light":
		return core.ThemeLight
	case "Dark":
		return core.ThemeDark
	}
	return core.ThemeAuto
}
