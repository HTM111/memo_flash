package ui

import (
	"fmt"
	"strconv"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/tree"
)

type ExploreView struct {
	core.Frame

	App          *App
	searchQuery  string
	contentFrame *core.Frame
}

func (ev *ExploreView) Init() {
	ev.Frame.Init()
	ev.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Grow.Set(1, 1)
	})
	tree.AddChild(ev, func(w *core.TextField) {
		w.SetType(core.TextFieldOutlined)
		w.Styler(func(s *styles.Style) {
			s.Grow.Set(1, 0)
			s.Max.Zero()
			s.Justify.Items = styles.Center
		})
		w.SetLeadingIcon(icons.Search)
		w.SetPlaceholder("Search cards...")
		w.OnInput(func(e events.Event) {
			ev.searchQuery = w.Text()
			ev.contentFrame.Update()
		})

	})
	tree.AddChild(ev, func(w *core.Frame) {
		ev.contentFrame = w

		w.Styler(func(s *styles.Style) {
			s.Grow.Set(1, 0)
			s.Direction = styles.Column
		})
		w.Maker(func(p *tree.Plan) {
			ev.makeContent(p)
		})

	})
}
func (ev *ExploreView) makeContent(p *tree.Plan) {
	searchResults := ev.App.StudyManager.SearchCards(ev.searchQuery)
	if len(searchResults) == 0 {
		tree.AddAt(p, "empty-state", func(w *emptyState) {
			w.Updater(func() {
				if ev.searchQuery == "" {
					w.SetMessage("No cards in deck")
				} else {
					w.SetMessage(fmt.Sprintf("'%s' not found", ev.searchQuery))
				}
			})
			w.SetIcon(icons.Search)
		})
	} else {
		for _, card := range searchResults {
			tree.AddAt(p, strconv.Itoa(card.ID), func(w *Card) {
				w.Updater(func() {
					w.SetData(card)
				})
				w.SetEdit(func() {
					ShowCardDialog(ev.App, &CardData{
						Front: card.Front,
						Back:  card.Back,
					}, true, func(cd *CardData) {
						err := ev.App.StudyManager.EditCard(card, cd.Front, cd.Back)
						if err != nil {
							core.ErrorDialog(ev.App, err, "Can't edit card")
							return
						}
						w.Update()
					})
				})
				w.SetDelete(func() {
					if err := ev.App.StudyManager.DeleteCard(card); err != nil {
						core.ErrorSnackbar(ev, err, "Error While deleting card")
					}
					ev.App.UpdateItemInDeck(card.ParentDeckId)
					ev.contentFrame.Update()
				})
			})
		}
	}
}
func (app *App) UpdateItemInDeck(id int) {
	decklist := app.Tabs.TabByName("Decks").ChildByName("deck-list", 0).(*core.Frame)
	if decklist != nil {
		decklist.ChildByName(strconv.Itoa(id), 0).(*Deck).Update()
	}
}
