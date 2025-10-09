package ui

import (
	"image/color"
	"log"
	"memoflash/internal/models"
	"memoflash/internal/services"
	"memoflash/internal/values"
	"memoflash/pkg/fsrs"
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

type DeckTab struct {
	core.Frame
	services *services.Service
	Decks    []*models.Deck
	deckList *core.Frame
	decksMap map[int]*models.Deck
}

func (dt *DeckTab) Init() {
	dt.decksMap = make(map[int]*models.Deck)
	dt.Frame.Init()
	dt.Styler(func(s *styles.Style) {
		s.Margin.SetAll(units.Dp(15))
		s.Grow.Set(1, 1)
		s.Direction = styles.Column
	})
	dt.createDeckTab()
	dt.Updater(func() {
		if dt.Decks == nil {
			list, err := dt.services.GetDecks()
			if err != nil {
				log.Println("Err", err.Error())
			}
			if len(list) == 0 {
				list = []*models.Deck{}
			}
			dt.Decks = list
			for _, item := range list {
				dt.decksMap[item.ID] = item
			}
		}
	})

}

func (dt *DeckTab) createDeckTab() {
	dt.makeDeckHeader(dt)
	dt.makeDeckListContainer(dt)

}

func (dt *DeckTab) makeDeckListContainer(frame tree.Node) {

	tree.AddChildAt(frame, "deck-list", func(w *core.Frame) {
		w.Styler(func(s *styles.Style) {
			s.Direction = styles.Column
			s.Grow.Set(1, 1)
			s.Gap.Set(units.Dp(12))
		})
		dt.deckList = w
		w.Maker(func(p *tree.Plan) {
			if len(dt.Decks) > 0 {
				dt.makeDeckList(p, dt.Decks)
			} else {
				EmptyState(p, "No Decks Found", icons.Info)
			}
		})

	})

}

func (dt *DeckTab) makeDeckHeader(parent tree.Node) {
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
				dueCards, err := dt.services.GetAllDueCards()
				if err != nil {
					core.ErrorSnackbar(dt, err, "Error Getting Due Cards")
					return
				}
				if len(dueCards) == 0 {
					core.MessageDialog(dt, "No Due cards to study")
					return
				}
				dt.HandleStudy(dueCards)
			})
		})
		tree.AddChildAt(w, "deck-create-button", func(w *core.Button) {
			w.SetIcon(icons.Add)
			w.SetText("Create Deck")
			w.Styler(func(s *styles.Style) {
				s.Padding.SetAll(units.Dp(12))
			})
			w.OnClick(func(e events.Event) {
				ShowDeckDialog(dt, &DeckData{},
					false, func(dd *DeckData) {
						id, err := dt.services.CreateDeck(dd.Title, dd.Description, dd.CategoryColorIndex)
						if err != nil {
							core.ErrorSnackbar(dt, err, "Error Creating Deck")
							return
						}
						item := &models.Deck{
							ID:            id,
							Title:         dd.Title,
							Description:   dd.Description,
							CategoryIndex: dd.CategoryColorIndex,
						}
						dt.Decks = append(dt.Decks, item)
						dt.decksMap[id] = item
						dt.deckList.Update()

					})

			})
		})
	})
}
func (dt *DeckTab) handleActions(w *Deck, deck *models.Deck) {
	w.OnAddCard(func() {
		ShowCardDialog(dt, &CardData{}, false, func(card *CardData) {
			err := dt.services.CreateCard(card.Front, card.Back, deck.ID)
			if err != nil {
				core.ErrorSnackbar(dt, err, "Error Creating Card")
				return
			}
			deck.TotalCards += 1
			deck.DueCards += 1
			w.Update()

		})

	})
	w.OnEdit(func() {
		ShowDeckDialog(dt, &DeckData{
			Title:              deck.Title,
			Description:        deck.Description,
			CategoryColorIndex: deck.CategoryIndex,
		},
			true, func(dd *DeckData) {
				err := dt.services.EditDeck(deck.ID, dd.Title, dd.Description, dd.CategoryColorIndex)
				if err != nil {
					core.ErrorSnackbar(dt, err, "Error Updating Deck")
					return
				}
				deck.Title = dd.Title
				deck.CategoryIndex = dd.CategoryColorIndex
				deck.Description = dd.Description
				w.Update()
			})
	})
	w.OnExplore(func() {
		pm := core.NewBody()
		tree.AddChild(pm, func(w *ExploreView) {
			w.service = dt.services
			w.deckListFrame = dt.deckList
			w.deck = deck
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
			err := dt.services.DeleteDeck(deck.ID)
			if err != nil {
				core.ErrorSnackbar(dt, err, "Error Deleting Deck")
				return
			}
			for i, deckitem := range dt.Decks {
				if deckitem.ID == deck.ID {
					dt.Decks = append(dt.Decks[:i], dt.Decks[i+1:]...)
				}
			}
			delete(dt.decksMap, deck.ID)
			dt.UpdateList()
		}
		if deck.TotalCards == 0 {
			deletAction()
			return
		}
		WarningDialog(dt, "Delete Deck ? ", "this action is irreversible", "Delete", func() {
			deletAction()
		})
	})

	w.OnStudy(func() {
		if deck.TotalCards == 0 {
			core.MessageDialog(dt, "No cards to study")
			return
		}
		dueCards, err := dt.services.GetDueCardsFromDeck(deck.ID)
		if err != nil {
			core.ErrorSnackbar(dt, err, "Error Getting Due Cards")
			return
		}
		if len(dueCards) == 0 {
			core.MessageDialog(dt, "No cards to study")
			return
		}

		dt.HandleStudy(dueCards)

	})

}
func (dt *DeckTab) HandleStudy(dueCards []*models.Card) {
	d := core.NewBody("Back to Decks")
	var previousDeckId = dueCards[0].ParentDeckId
	var isSame bool = true
	pages := core.NewPages(d)
	pages.AddPage("main", func(pg *core.Pages) {
		p := core.NewFrame(pg)
		p.Styler(func(s *styles.Style) {
			s.Grow.Set(1, 1)
			s.CenterAll()
		})

		tree.AddChild(p, func(w *StudyPage) {
			w.Cards = dueCards
			w.OnEach = func(card *models.Card, rating values.Difficulty) error {
				dt.decksMap[card.ParentDeckId].DueCards--
				fsrs.Review(rating, card)
				if previousDeckId != card.ParentDeckId {
					isSame = false
				}
				previousDeckId = card.ParentDeckId
				return dt.services.UpdateInterval(card.ID,
					card.Interval,
					card.Difficulty,
					card.Stability)
			}

			w.OnDone = func() { pages.Open("status-page") }
		})
	})

	pages.AddPage("status-page", func(pg *core.Pages) {
		fr := core.NewFrame(pg)
		fr.Styler(func(s *styles.Style) {
			s.CenterAll()
			s.Grow.Set(1, 1)
		})
		StatePage(fr)
	})

	d.OnClose(func(e events.Event) {
		if isSame {
			dt.services.UpdateReadTime(previousDeckId)
		}
		dt.UpdateList()
	})

	d.RunFullDialog(dt)
}
func (dt *DeckTab) UpdateItemInList(id int) {
	item := dt.deckList.ChildByName(strconv.Itoa(id), 0)
	if item != nil {
		item.(*Deck).Update()
	}
}
func (dt *DeckTab) UpdateList() {
	dt.deckList.Update()
}
func (dt *DeckTab) makeDeckList(p *tree.Plan, items []*models.Deck) {
	for _, deck := range items {
		tree.AddAt(p, strconv.Itoa(deck.ID), func(w *Deck) {

			dt.handleActions(w, deck)
			w.Updater(func() {
				w.setData(deck)
			})

		})
	}
}
