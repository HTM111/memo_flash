package ui

import (
	"memoflash/internal/models"
	"memoflash/internal/services"
	"path/filepath"
	"slices"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/states"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/tree"
)

var Settings = &AppSettings{
	SettingsBase: core.SettingsBase{
		Name: "memoflash",
		File: filepath.Join(core.TheApp.DataDir(), "memoflash", "settings.toml"),
	},
}

func init() {
	core.TheApp.SetName("memoflash")
	core.AllSettings = slices.Insert(core.AllSettings, 1, core.Settings(Settings))
	core.TheApp.SetSceneInit(func(sc *core.Scene) {
		sc.Styler(func(s *styles.Style) {
			s.Background = colors.Scheme.Surface
			s.Padding.Set(units.Dp(10))
			s.Margin.Zero()
			s.Grow.Set(1, 1)
		})
		sc.SetWidgetInit(func(w core.Widget) {
			switch w := w.(type) {
			case *core.TextField:
				w.SetType(core.TextFieldOutlined)
				w.Styler(func(s *styles.Style) {
					s.Border.Radius = styles.BorderRadiusExtraSmall
					s.Background = colors.Scheme.SurfaceContainer
					s.Border.Width.Zero()
					s.Border.Color.Zero()

					if s.Is(states.Focused) {
						s.Border.Width.Set(units.Dp(2))
						s.Border.Color.Set(colors.Scheme.Primary.Base)
					}
				})
			case *core.Button:
				w.Styler(func(s *styles.Style) {
					s.Border.Radius = styles.BorderRadiusSmall
				})
			}
		})
	})
}

type App struct {
	core.Frame
	Services *services.Service
	Settings *AppSettings
	Tabs     *core.Tabs
	Decks    []*models.Deck
	DeckMap  map[int]*models.Deck
}

func NewMemoFlashWindow(service *services.Service) {
	appName := "MemoFlash"
	b := core.NewBody(appName).SetTitle(appName)
	app := tree.New[App](b)
	app.Services = service
	app.CreateApp()
	b.RunMainWindow()
}

func (app *App) Init() {
	app.DeckMap = make(map[int]*models.Deck, 0)
	app.Frame.Init()

}

func (app *App) CreateApp() {
	app.FetchDecks()
	app.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Padding.Zero()
		s.Grow.Set(1, 1)
	})
	tree.AddChildAt(app, "main-content", func(w *core.Frame) {
		w.Styler(func(s *styles.Style) {
			s.Direction = styles.Column
			s.Grow.Set(1, 1)
		})
		app.createMainContent(w)
	})
}
func (app *App) createMainContent(parent *core.Frame) {
	tree.AddChildAt(parent, "tabs", func(w *core.Tabs) {
		w.SetType(core.NavigationDrawer)
		w.Styler(func(s *styles.Style) {
			s.Grow.Set(1, 1)
		})
		app.createTabPanels(w)
	})
}

func (app *App) createTabPanels(tabs *core.Tabs) {
	app.Tabs = tabs

	frameStudy, tabStudy := tabs.NewTab("Study")
	tabStudy.SetIcon(icons.School)
	tree.AddChildAt(frameStudy, "decks-section", func(studyTab *StudyTab) {
		studyTab.deckrepo = app
		studyTab.services = app.Services

	})

	frameDeck, decktab := tabs.NewTab("Decks")
	decktab.SetIcon(icons.List)
	tree.AddChildAt(frameDeck, "decks-section", func(deckTab *DeckTab) {
		deckTab.deckrepo = app
		deckTab.service = app.Services
	})
}
