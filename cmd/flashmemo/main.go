package main

import (
	"memoflash/internal/db"
	"memoflash/internal/services"
	"memoflash/internal/ui"
	"path"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
)

func main() {
	services, err := initialization()
	if err != nil {
		ShowError(err)
		return
	}

	ui.NewMemoFlashWindow(services)

}
func initialization() (*services.Service, error) {
	if err := core.LoadAllSettings(); err != nil {
		return nil, err
	}
	p := path.Join(core.TheApp.AppDataDir(), "memoflash.db")
	db, err := db.SetupDatabase(p)
	if err != nil {
		return nil, err
	}

	if err := db.InitSchema(); err != nil {
		return nil, err
	}

	service := &services.Service{
		CardService: services.NewCardService(db),
		DeckService: services.NewDeckService(db),
	}
	if err != nil {
		return nil, err
	}
	return service, nil

}
func ShowError(err error) {

	errorBody := core.NewBody("MemoFlash - Error").SetTitle("MemoFlash - Error")
	errorBody.Styler(func(s *styles.Style) {
		s.Min.Set(units.Dp(350), units.Dp(200))
	})

	core.NewText(errorBody).SetType(core.TextBodyMedium).SetText("Cannot create Study Manager: " + err.Error())
	errorBody.AddBottomBar(func(bar *core.Frame) {
		core.NewButton(bar).SetText("Exit").OnClick(func(e events.Event) {
			errorBody.Scene.Close()
		})
	})
	win := errorBody.NewWindow()
	win.SetUseMinSize(true)
	win.RunMain()
}
