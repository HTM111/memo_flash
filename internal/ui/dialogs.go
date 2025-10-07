package ui

import (
	"strings"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/states"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/text/rich"
)

type CardData struct {
	Front    string
	Back     string
	KeepOpen bool
}
type DeckData struct {
	Title              string
	Description        string
	CategoryColorIndex int
}

func ShowCardDialog(ctx core.Widget, data *CardData, isEdit bool, onAccept func(*CardData)) {
	title := "Create a new flashcard"
	if isEdit {
		title = "Edit your flashcard"
	}
	isDisabled := !isEdit
	d := core.NewBody(title)
	core.NewText(d).SetType(core.TextBodyMedium).SetText(title)

	core.NewText(d).SetText("Front").Styler(func(s *styles.Style) {
		s.Font.Weight = rich.Bold
	})

	frontField := core.NewTextField(d).SetPlaceholder("Enter front text")
	frontField.SetType(core.TextFieldOutlined)
	frontField.Styler(func(s *styles.Style) {
		s.Grow.Set(1, 0)
		s.Max.Zero()
	})
	frontField.SetText(data.Front)
	frontField.OnChange(func(e events.Event) {
		data.Front = frontField.Text()
	})

	core.NewText(d).SetText("Back").Styler(func(s *styles.Style) {
		s.Font.Weight = rich.Bold
	})

	backField := core.NewTextField(d).SetPlaceholder("Enter back text")
	backField.Styler(func(s *styles.Style) {
		s.Grow.Set(1, 0)
		s.Max.Zero()
	})
	backField.SetText(data.Back)
	backField.OnChange(func(e events.Event) {
		data.Back = backField.Text()
	})

	if !isEdit {
		var keepOpenSwitch *core.Switch

		keepOpenSwitch = core.NewSwitch(d).SetText("Keep Open")
		keepOpenSwitch.SetChecked(data.KeepOpen)
		keepOpenSwitch.OnChange(func(e events.Event) {
			data.KeepOpen = keepOpenSwitch.IsChecked()
		})
	}

	d.AddBottomBar(func(bar *core.Frame) {
		d.AddCancel(bar)
		create := core.NewButton(bar)
		if isEdit {
			create.SetText("Save")
		} else {
			create.SetText("Create")
		}
		updateButton := func() {
			isDisabled = len(strings.TrimSpace(frontField.Text())) == 0 || len(strings.TrimSpace(backField.Text())) == 0
			create.Update()
		}

		create.Updater(func() {
			create.SetState(isDisabled, states.Disabled)
		})
		frontField.OnInput(func(e events.Event) {
			updateButton()
		})
		backField.OnInput(func(e events.Event) {
			updateButton()
		})

		if isEdit {
			create.SetText("Save")
		} else {
			create.SetText("Create")
		}
		create.OnClick(func(e events.Event) {
			if onAccept != nil {
				if data.KeepOpen {
					backField.SetText("")
					frontField.SetText("")
					isDisabled = true
					create.Update()
				} else {
					d.Close()
				}
				onAccept(data)
			}

		})
	})

	dialog := d.NewDialog(ctx)
	dialog.SetDisplayTitle(true)
	dialog.SetResizable(false)
	dialog.Run()
}

func ShowDeckDialog(ctx core.Widget, data *DeckData, isEdit bool, onAccept func(*DeckData)) {
	isDisabled := !isEdit
	title := "Create a new deck"
	if isEdit {
		title = "Edit your deck"
	}
	d := core.NewBody(title)

	core.NewText(d).SetType(core.TextBodyMedium).SetText(title)
	core.NewText(d).SetText("Title")
	titleField := core.NewTextField(d).
		SetPlaceholder("Enter Title").
		SetType(core.TextFieldOutlined)
	titleField.Styler(func(s *styles.Style) {
		s.Grow.Set(1, 0)
		s.Max.Zero()

	})
	titleField.SetText(data.Title)
	titleField.OnChange(func(e events.Event) {
		data.Title = titleField.Text()
	})
	core.NewText(d).SetText("Description")

	descField := core.NewTextField(d).SetPlaceholder("Enter Description")
	descField.SetType(core.TextFieldOutlined)

	descField.Styler(func(s *styles.Style) {
		s.Grow.Set(1, 0)
		s.Max.Zero()
	})
	descField.SetText(data.Description)
	descField.OnChange(func(e events.Event) {
		data.Description = descField.Text()
	})
	core.NewText(d).SetText("Category")

	categoryFrame := core.NewFrame(d)
	categoryFrame.Styler(func(s *styles.Style) {
		s.Border.Width.Set(units.Zero())
		s.Grow.Set(1, 0)
		s.Gap.Zero()
		s.Padding.Zero()
	})
	for i, color := range CategoryColors {
		btn := core.NewButton(categoryFrame)
		btn.Styler(func(s *styles.Style) {
			s.Background = colors.Uniform(color)
			s.Border.Radius = styles.BorderRadiusExtraSmall
			s.Min.Set(units.Dp(33), units.Dp(33))
			s.Grow.Set(1, 0)
			s.Padding.Zero()

			if s.Is(states.Selected) {
				s.Border.Width.Set(units.Dp(2))
				s.Border.Color.Set(colors.Uniform(colors.White))
			}
		})

		index := i
		btn.SetSelected(data.CategoryColorIndex == index)
		btn.OnClick(func(e events.Event) {
			data.CategoryColorIndex = index
			for j, child := range categoryFrame.Children {
				if childBtn, ok := child.(*core.Button); ok {
					childBtn.SetSelected(j == index)
				}
			}
		})
	}

	d.AddBottomBar(func(bar *core.Frame) {
		d.AddCancel(bar)
		create := d.AddOK(bar)
		updateButton := func() {
			isDisabled = len(strings.TrimSpace(titleField.Text())) == 0
			create.Update()
		}
		create.Updater(func() {

			create.SetState(isDisabled, states.Disabled)
		})
		titleField.OnInput(func(e events.Event) {
			updateButton()
		})

		if isEdit {
			create.SetText("Save")
		} else {
			create.SetText("Create")
		}
		create.OnClick(func(e events.Event) {
			if onAccept != nil {
				d.Close()
				onAccept(data)
			}
		})
	})
	dialog := d.NewDialog(ctx)
	dialog.SetDisplayTitle(true)
	dialog.SetResizable(false)
	dialog.Run()

}
func WarningDialog(ctx core.Widget, title, message string, buttonText string, onYes func()) {
	dialog := core.NewBody(title)
	core.NewText(dialog).SetText(message)
	dialog.AddBottomBar(func(bar *core.Frame) {
		dialog.AddCancel(bar)
		btn := dialog.AddOK(bar)
		btn.SetText(buttonText)
		btn.Styler(func(s *styles.Style) {
			s.Color = colors.Uniform(colors.White)
			s.Background = colors.Uniform(colors.Red)
		})
		btn.OnClick(func(e events.Event) {
			if onYes != nil {
				dialog.Close()
				onYes()
			}
		})
	})
	d := dialog.NewDialog(ctx)
	d.SetResizable(false)
	d.Run()
}
