package ui

import (
	"fmt"
	"memoflash/internal/models"
	"memoflash/internal/services"
	"strconv"
	"strings"

	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/tree"
)

type ExploreView struct {
	core.Frame
	deck          *models.Deck
	service       *services.Service
	searchQuery   string
	Cards         []*models.Card
	deckListFrame *core.Frame
	contentFrame  *core.Frame
}

func (ev *ExploreView) Init() {
	ev.Frame.Init()

	ev.Styler(func(s *styles.Style) {
		s.Direction = styles.Column
		s.Grow.Set(1, 1)
	})
	ev.Updater(func() {
		if ev.Cards == nil {
			cards, err := ev.service.GetCardsByDeck(ev.deck.ID)
			if err != nil {
				core.ErrorDialog(ev, err, "Error Getting Deck")
				ev.Cards = []*models.Card{}
				return
			}
			ev.Cards = cards
		}
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
	searchResults := ev.SearchCards(ev.searchQuery)
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
					ShowCardDialog(ev, &CardData{
						Front: card.Front,
						Back:  card.Back,
					}, true, func(cd *CardData) {
						err := ev.service.EditCard(card.ID, cd.Front, cd.Back)
						if err != nil {
							core.ErrorDialog(ev, err, "Can't edit card")
							return
						}
						card.Front = cd.Front
						card.Back = cd.Front
						w.Update()
					})
				})
				w.SetDelete(func() {
					if err := ev.service.DeleteCard(card.ID); err != nil {
						core.ErrorSnackbar(ev, err, "Error While deleting card")
						return
					}
					for index, cardItem := range ev.Cards {
						if cardItem.ID == card.ID {
							ev.Cards = append(ev.Cards[:index], ev.Cards[index+1:]...)
						}
					}
					ev.deck.TotalCards--
					if card.IsDue() {
						ev.deck.DueCards--
					}
					ev.deckListFrame.Update()
					ev.contentFrame.Update()
				})
			})
		}
	}
}
func (ev *ExploreView) SearchCards(query string) []*models.Card {

	if query == "" {
		return ev.Cards
	}
	var cards []*models.Card
	queryLower := strings.ToLower(query)
	for _, card := range ev.Cards {
		if card == nil {
			continue
		}
		if strings.Contains(strings.ToLower(card.Front), queryLower) ||
			strings.Contains(strings.ToLower(card.Back), queryLower) {
			cards = append(cards, card)
		}
	}
	return cards
}
