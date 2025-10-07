package ui

import (
	"fmt"
	"image/color"
	"memoflash/internal/models"
	"memoflash/internal/values"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/abilities"
	"cogentcore.org/core/styles/states"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/text/rich"
	"cogentcore.org/core/tree"
)

type StudyPage struct {
	core.Pages
	Cards            []*models.Card
	Title            string
	CurrentCardIndex int
	ShowFront        bool
	showButtons      bool
	OnEach           func(card *models.Card, rating values.Difficulty) error
	OnDone           func()
}

func (sd *StudyPage) Init() {
	sd.Pages.Init()
	sd.ShowFront = true
	sd.showButtons = false
	sd.CurrentCardIndex = 0

}
func NewStudyPage(parent tree.Node) *StudyPage {
	sd := tree.New[StudyPage](parent)
	return sd
}
func (sd *StudyPage) Create(dueCards []*models.Card, title string, onDone func(), onEach func(card *models.Card, rating values.Difficulty) error) {
	sd.Cards = dueCards
	sd.Title = title
	sd.OnEach = onEach
	sd.OnDone = onDone
	sd.AddPage("study", sd.makeStudyPage)
	sd.AddPage("congratulations", sd.makeCongratulationsPage)
}
func (sd *StudyPage) handleRating(rating values.Difficulty) {
	if sd.OnEach != nil && len(sd.Cards) > sd.CurrentCardIndex {
		err := sd.OnEach(sd.Cards[sd.CurrentCardIndex], rating)
		if err != nil {
			core.ErrorSnackbar(sd, err, "Error Updating Interval")
			return
		}
	}

	sd.CurrentCardIndex++
	if sd.CurrentCardIndex < len(sd.Cards) {
		sd.ShowFront = true
		sd.showButtons = false
		sd.Update()
	} else {
		if sd.OnDone != nil {
			sd.OnDone()
		}
		sd.Open("congratulations")
	}
}
func (sd *StudyPage) makeStudyPage(pg *core.Pages) {
	fr := core.NewFrame(pg)
	fr.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Grow.Set(1, 1)
		s.CenterAll()
	})

	container := core.NewFrame(fr)
	container.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Border.Radius.SetAll(units.Dp(10))
		s.Padding.SetAll(units.Dp(40))
		s.Background = colors.Scheme.SurfaceContainer
	})

	progressFrame := core.NewFrame(container)
	progressFrame.Styler(func(s *styles.Style) {
		s.Direction = styles.Row
		s.Grow.Set(1, 0)
		s.Align.Items = styles.Center
		s.Background = colors.Scheme.SurfaceContainer
	})

	textfr := core.NewText(progressFrame)
	textfr.Updater(func() {
		textfr.SetText(fmt.Sprintf("%d/%d", sd.CurrentCardIndex+1, len(sd.Cards)))
	})

	meter := core.NewMeter(progressFrame)
	meter.Styler(func(s *styles.Style) {
		s.Grow.Set(1, 0)
		s.Min.X.Zero()
		s.Max.X.Zero()
	})
	meter.Updater(func() {
		meter.SetMax(float32(len(sd.Cards)))
		meter.SetValue(float32(sd.CurrentCardIndex + 1))
	})

	cardFrame := core.NewFrame(container)
	cardFrame.Styler(func(s *styles.Style) {
		s.Margin.SetTop(units.Dp(10))
		s.SetAbilities(true, abilities.Clickable)
		s.Background = colors.Scheme.Surface
		s.Border.Radius.SetAll(units.Dp(10))
		s.Grow.Set(1, 1)
		s.Min.Set(units.Dp(450))
		s.Max.Set(units.Dp(700))
		s.Direction = styles.Column
	})
	cardFrame.OnClick(func(e events.Event) {
		sd.ShowFront = !sd.ShowFront
		sd.showButtons = true
		container.Update()
	})

	mainContent := core.NewFrame(cardFrame)
	mainContent.Styler(func(s *styles.Style) {
		s.Grow.Set(1, 1)
		s.CenterAll()
	})

	titleText := core.NewText(mainContent)
	titleText.SetType(core.TextHeadlineLarge)
	titleText.Styler(func(s *styles.Style) {
		s.SetNonSelectable()
		s.Font.Weight = rich.Bold
	})
	titleText.OnClick(func(e events.Event) {
		cardFrame.Send(events.Click, e)
	})
	titleText.Updater(func() {
		if sd.ShowFront {
			titleText.SetText(sd.Cards[sd.CurrentCardIndex].Front)
		} else {
			titleText.SetText(sd.Cards[sd.CurrentCardIndex].Back)
		}
	})

	bottomRow := core.NewFrame(cardFrame)
	bottomRow.Styler(func(s *styles.Style) {
		s.Direction = styles.Row
		s.Justify.Content = styles.End
		s.Grow.Set(1, 0)
	})

	cornerText := core.NewText(bottomRow)
	cornerText.SetType(core.TextBodySmall)
	cornerText.Styler(func(s *styles.Style) {
		s.Padding.SetAll(units.Dp(10))
		s.SetNonSelectable()
	})
	cornerText.SetReadOnly(true)
	cornerText.Updater(func() {
		if sd.ShowFront {
			cornerText.SetText("Tap to Reveal")
		} else {
			cornerText.SetText("Tap to Flip Back")
		}
	})
	cornerText.OnClick(func(e events.Event) {
		cardFrame.Send(events.Click, e)
	})
	buttonFrame := core.NewFrame(container)

	buttonFrame.Styler(func(s *styles.Style) {
		s.Justify.Content = styles.Center
		s.Gap.Set(units.Dp(5))
		s.Grow.Set(1, 0)
		s.Margin.Top = units.Dp(20)
	})
	buttonFrame.Updater(func() {
		buttonFrame.SetState(!sd.showButtons, states.Invisible)
	})

	easyBtn := core.NewButton(buttonFrame)
	easyBtn.Styler(func(s *styles.Style) {
		s.Background = colors.Uniform(color.RGBA{160, 210, 160, 255}) // Soft green
		s.Color = colors.Uniform(color.RGBA{0, 60, 0, 255})
		s.Font.Weight = rich.Bold
		s.BoxShadow = []styles.Shadow{
			{
				OffsetX: units.Zero(),
				OffsetY: units.Dp(3),
				Spread:  units.Zero(),
				Color:   colors.Uniform(color.RGBA{0, 60, 0, 255}),
			},
		}
	})
	easyBtn.SetText("Easy").OnClick(func(e events.Event) {
		sd.handleRating(values.Easy)
	})

	goodBtn := core.NewButton(buttonFrame)
	goodBtn.Styler(func(s *styles.Style) {
		bg := color.RGBA{100, 180, 170, 255}
		fg := color.RGBA{0, 70, 60, 255}

		s.Background = colors.Uniform(bg)
		s.Color = colors.Uniform(fg)
		s.Font.Weight = rich.Bold
		s.BoxShadow = []styles.Shadow{{
			OffsetX: units.Zero(),
			OffsetY: units.Dp(3),
			Spread:  units.Zero(),
			Color:   colors.Uniform(fg),
		}}
	})
	goodBtn.SetText("Good").OnClick(func(e events.Event) {
		sd.handleRating(values.Good)
	})

	hardBtn := core.NewButton(buttonFrame)
	hardBtn.Styler(func(s *styles.Style) {
		s.Background = colors.Uniform(color.RGBA{255, 200, 140, 255}) // Soft orange
		s.Color = colors.Uniform(color.RGBA{100, 40, 0, 255})
		s.Font.Weight = rich.Bold
		s.BoxShadow = []styles.Shadow{
			{
				OffsetX: units.Zero(),
				OffsetY: units.Dp(3),
				Spread:  units.Zero(),
				Color:   colors.Uniform(color.RGBA{100, 40, 0, 255}),
			},
		}
	})
	hardBtn.SetText("Hard").OnClick(func(e events.Event) {
		sd.handleRating(values.Hard)
	})

	againBtn := core.NewButton(buttonFrame)
	againBtn.Styler(func(s *styles.Style) {
		bgColor := color.RGBA{240, 140, 140, 255}
		textColor := color.RGBA{90, 0, 0, 255}
		s.Background = colors.Uniform(bgColor)
		s.Color = colors.Uniform(textColor)
		s.Font.Weight = rich.Bold
		s.BoxShadow = []styles.Shadow{
			{
				OffsetX: units.Zero(),
				OffsetY: units.Dp(3),
				Spread:  units.Zero(),
				Color:   colors.Uniform(textColor),
			},
		}
	})
	againBtn.SetText("Again").OnClick(func(e events.Event) {
		sd.handleRating(values.Again)
	})
}

func (sd *StudyPage) makeCongratulationsPage(pg *core.Pages) {
	r := core.NewFrame(pg)
	pg.Styler(func(s *styles.Style) {
		s.CenterAll()
	})
	r.Styler(func(s *styles.Style) {
		s.Gap.Set(units.Dp(20))
		s.Direction = styles.Column
		s.Border.Radius.SetAll(units.Dp(10))
		s.CenterAll()

		s.Padding.Set(units.Dp(50))
		s.Background = colors.Scheme.SurfaceContainerLow
	})
	core.NewText(r).SetType(core.TextHeadlineSmall).
		SetText("🎉 Congratulations! 🎉").
		Styler(func(s *styles.Style) {
			s.SetTextWrap(false)
		})
	txt := core.NewText(r).SetType(core.TextHeadlineSmall)
	txt.SetText("You've completed Your Session!")
	txt.Styler(func(s *styles.Style) {
		s.SetTextWrap(false)
	})
}
