package ui

import (
	"image/color"
	"memoflash/internal/models"
	"memoflash/internal/values"
	"strconv"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/text/rich"
	"cogentcore.org/core/tree"
)

var CategoryColors = []color.Color{
	color.RGBA{102, 187, 106, 255},
	color.RGBA{66, 165, 245, 255},
	color.RGBA{239, 83, 80, 255},
	color.RGBA{255, 167, 38, 255},
	color.RGBA{171, 71, 188, 255},
	color.RGBA{38, 198, 218, 255},
}

func (app *App) createDeckTab(frame *core.Frame) {
	frame.Styler(func(s *styles.Style) {
		s.Margin.SetAll(units.Dp(15))

	})

	app.makeDeckHeader(frame)
	app.makeDeckListContainer(frame)

}

func (app *App) makeDeckListContainer(frame *core.Frame) {
	tree.AddChildAt(frame, "deck-list", func(w *core.Frame) {
		w.Styler(func(s *styles.Style) {
			s.Direction = styles.Column
			s.Grow.Set(1, 1)
			s.Gap.Set(units.Dp(12))
		})

		w.Maker(func(p *tree.Plan) {
			list := app.StudyManager.GetDecks()
			if len(list) > 0 {
				app.makeDeckList(p, list)
			} else {
				app.emptyState(p, "No Decks Found", icons.Info)
			}
		})

	})

}

func (app *App) makeDeckHeader(parent *core.Frame) {
	tree.AddChildAt(parent, "header", func(w *core.Frame) {
		w.Styler(func(s *styles.Style) {
			s.Grow.Set(1, 0)
			s.Align.Items = styles.Center
			s.Padding.Bottom = units.Dp(5)
		})

		tree.AddChildAt(w, "deck-title", func(w *core.Frame) {
			w.Styler(func(s *styles.Style) {
				s.Grow.Set(1, 0)
			})
			tree.AddChild(w, func(w *core.Text) {
				w.SetText("Your Decks")

				w.SetType(core.TextHeadlineMedium)
				w.Styler(func(s *styles.Style) {
					s.Font.Weight = rich.Bold
				})
			})
		})
		tree.AddChildAt(w, "deck-quick-study", func(w *core.Button) {
			w.SetIcon(icons.PlayArrow)
			w.SetText("Study Due Cards")
			w.Styler(func(s *styles.Style) {
				s.Padding.SetAll(units.Dp(12))
			})
			w.OnClick(func(e events.Event) {
				if app.StudyManager.GetDueCount() == 0 {
					core.MessageDialog(app, "No Due Cards", "Study")
					return
				}
				dueCards, err := app.StudyManager.GetAllDueCards()
				if err != nil {
					core.ErrorSnackbar(app, err, "Error Getting Due Cards")
					return
				}

				fullDialog := core.NewBody("Back to Decks")

				page := NewStudyPage(fullDialog)
				page.Create(dueCards, "Due Cards", nil, func(card *models.Card, rating values.Difficulty) error {
					return app.StudyManager.UpdateInterval(card, rating)
				})
				fullDialog.OnClose(func(e events.Event) {
					app.updateDeckList()
				})
				fullDialog.RunFullDialog(app)
			})
		})
		tree.AddChildAt(w, "deck-create-button", func(w *core.Button) {
			w.SetIcon(icons.Add)
			w.SetText("Create Deck")
			w.Styler(func(s *styles.Style) {
				s.Padding.SetAll(units.Dp(12))
			})
			w.OnClick(func(e events.Event) {
				ShowDeckDialog(app, &DeckData{},
					false, func(dd *DeckData) {
						err := app.StudyManager.CreateDeck(dd.Title, dd.Description, dd.CategoryColorIndex)
						if err != nil {
							core.ErrorSnackbar(app, err, "Error Creating Deck")
							return
						}
						parent.Update()
					})

			})
		})
	})
}
func (app *App) handleActions(w *Deck, deck *models.Deck) {
	w.OnAddCard(func() {
		ShowCardDialog(app, &CardData{}, false, func(card *CardData) {
			err := app.StudyManager.CreateCard(card.Front, card.Back, deck)
			if err != nil {
				core.ErrorSnackbar(app, err, "Error Creating Card")
				return
			}
			w.Update()

		})

	})
	w.OnEdit(func() {
		ShowDeckDialog(app, &DeckData{
			Title:              deck.Title,
			Description:        deck.Description,
			CategoryColorIndex: deck.CategoryIndex,
		},
			true, func(dd *DeckData) {
				err := app.StudyManager.UpdateDeck(deck, dd.Title, dd.Description, dd.CategoryColorIndex)
				if err != nil {
					core.ErrorSnackbar(app, err, "Error Updating Deck")
					return
				}
				w.Update()
			})
	})
	w.OnExplore(func() {
		app.StudyManager.SelectDeck(deck)
		pm := core.NewBody()
		tree.AddChild(pm, func(w *ExploreView) {
			w.App = app
		})
		pm.AddTopBar(func(bar *core.Frame) {
			closeBtn := core.NewButton(bar).SetIcon(icons.Close)
			bar.Styler(func(s *styles.Style) {
				s.Justify.Content = styles.End
				s.Align.Items = styles.Center
				s.Padding.SetAll(units.Dp(8))
			})
			closeBtn.SetType(core.ButtonAction)
			closeBtn.Styler(func(s *styles.Style) {
				s.Padding.SetAll(units.Dp(8))
			})
			closeBtn.OnClick(func(e events.Event) {
				app.StudyManager.ClearSelectedCards()
				pm.Close()
			})
		})

		dialog := pm.NewDialog(nil)
		dialog.SetModal(true)
		dialog.SetDisplayTitle(false)
		dialog.SetUseMinSize(true)
		dialog.SetResizable(false)
		dialog.Run()
	})

	w.OnDelete(func() {
		deletAction := func() {
			err := app.StudyManager.DeleteDeck(deck)
			if err != nil {
				core.ErrorSnackbar(app, err, "Error Deleting Deck")
				return
			}
			app.updateDeckList()
		}
		if deck.TotalCards == 0 {
			deletAction()
			return
		}
		WarningDialog(app, "Delete Deck ? ", "this action is irreversible", "Delete", func() {
			deletAction()
		})
	})

	w.OnStudy(func() {
		if deck.TotalCards == 0 {
			core.MessageDialog(app, "No cards to study")
			return
		}
		dueCards, err := app.StudyManager.GetDueCardsFromDeck(deck)
		if err != nil {
			core.ErrorSnackbar(app, err, "Error Getting Due Cards")
			return
		}
		if len(dueCards) == 0 {
			core.MessageDialog(app, "No cards to study")
			return
		}
		d := core.NewBody("Back to Decks")

		page := NewStudyPage(d)
		page.Create(dueCards, deck.Title, func() {
			app.StudyManager.UpdateReadTime(deck)
		}, func(card *models.Card, rating values.Difficulty) error {
			return app.StudyManager.UpdateInterval(card, rating)
		})

		d.OnClose(func(e events.Event) {
			app.StudyManager.UpdateReadTime(deck)
			app.StudyManager.RefreshDueCardsInDeck(deck)
			app.updateDeckList()
		})
		d.RunFullDialog(app)

	})

}
func (app *App) makeDeckList(p *tree.Plan, items []*models.Deck) {
	for _, deck := range items {
		tree.AddAt(p, strconv.Itoa(deck.ID), func(w *Deck) {
			app.handleActions(w, deck)
			w.Updater(func() {
				w.setData(deck)
			})

		})
	}
}
func (app *App) updateDeckList() {
	decklist := app.Tabs.TabByName("Decks").ChildByName("deck-list", 0)
	if decklist != nil {
		decklist.(*core.Frame).Update()
	}

}
