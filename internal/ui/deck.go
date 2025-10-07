package ui

import (
	"fmt"
	"memoflash/internal/models"
	"strconv"

	"cogentcore.org/core/colors"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/abilities"
	"cogentcore.org/core/styles/states"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/text/rich"
	"cogentcore.org/core/text/text"
	"cogentcore.org/core/tree"
)

type Deck struct {
	core.Frame
	deckdata  *models.Deck
	index     int
	onExplore func()
	onAddCard func()
	onEdit    func()
	onStudy   func()
	onMore    func()
	onDelete  func()
}

func (deck *Deck) Init() {
	deck.Frame.Init()

	deck.Styler(func(s *styles.Style) {
		s.CenterAll()
		s.Background = colors.Scheme.SurfaceContainerLow
		s.Border.Radius.Set(units.Dp(12))
		s.Padding.SetAll(units.Dp(12))
		s.Gap.Set(units.Dp(20))
		s.Grow.Set(1, 0)
		s.SetAbilities(true, abilities.Selectable)
	})

	tree.AddChild(deck, func(w *core.Frame) {
		w.Styler(func(s *styles.Style) {
			s.Border.Radius = styles.BorderRadiusFull
			s.Min.Set(units.Dp(15))
			s.Background = colors.Scheme.Primary.Base
		})
		w.Updater(func() {
			w.Styler(func(s *styles.Style) {
				s.Background = colors.Uniform(CategoryColors[deck.deckdata.CategoryIndex])
			})

		})
	})
	tree.AddChild(deck, func(w *core.Frame) {
		w.Styler(func(s *styles.Style) {
			s.Direction = styles.Column
			s.Grow.Set(1, 0)
		})

		tree.AddChild(w, func(w *core.Text) {
			w.Styler(func(s *styles.Style) {
				s.Font.Size.Set(20, units.UnitDp)
				s.Max.Y.Em(s.Text.LineHeight * 1)
				s.Font.Weight = rich.Bold
			})
			w.Updater(func() {
				if deck.deckdata != nil {
					w.SetText(deck.deckdata.Title)
				}
			})
		})

		tree.AddChild(w, func(w *core.Text) {
			w.Styler(func(s *styles.Style) {
				s.Color = colors.Scheme.OnSurfaceVariant

				s.Text.Align = text.Start
			})
			w.Updater(func() {
				if deck.deckdata != nil {
					w.SetText(deck.deckdata.Description)
				}
			})
		})

		tree.AddChild(w, func(w *core.Frame) {
			w.Styler(func(s *styles.Style) {
				s.Gap.Set(units.Dp(25))
				s.Grow.Set(1, 0)
			})
			tree.AddChild(w, func(w *core.Text) {
				w.Styler(func(s *styles.Style) {
					s.Color = colors.Scheme.OnSurfaceVariant

				})
				w.Updater(func() {
					if deck.deckdata != nil {
						w.SetText(fmt.Sprintf("Total Cards %d", deck.deckdata.TotalCards))
					}
				})
			})
			tree.AddChild(w, func(w *core.Text) {
				w.Styler(func(s *styles.Style) {
					s.Color = colors.Scheme.OnSurfaceVariant

				})
				w.Updater(func() {
					if deck.deckdata != nil {
						w.SetText(fmt.Sprintf("Due Cards %d", deck.deckdata.DueCards))
					}
				})
			})
		})

	})

	tree.AddChild(deck, func(w *core.Frame) {
		w.Styler(func(s *styles.Style) {
			s.Margin.SetTop(units.Dp(10))
			s.Margin.SetRight(units.Dp(10))
			s.Align.Items = styles.Center
		})
		deck.IconButton(w, icons.School, func(i *core.Icon) {
			i.OnClick(func(e events.Event) {
				if deck.onStudy != nil {
					deck.onStudy()
				}
			})

		})
		deck.IconButton(w, icons.Explore, func(i *core.Icon) {
			i.OnClick(func(e events.Event) {
				if deck.onExplore != nil {
					deck.onExplore()
				}
			})

		})
		deck.IconButton(w, icons.MoreHorizFill, func(i *core.Icon) {
			i.AddContextMenu(func(m *core.Scene) {
				core.NewButton(m).
					SetText("Delete").
					SetIcon(icons.Delete).
					OnClick(func(e events.Event) {
						if deck.onDelete != nil {
							deck.onDelete()
						}
					})
				core.NewButton(m).
					SetText("Add Card").
					SetIcon(icons.AddBoxFill).
					OnClick(func(e events.Event) {
						if deck.onAddCard != nil {
							deck.onAddCard()
						}
					})
				core.NewButton(m).
					SetText("Edit").
					SetIcon(icons.Edit).
					OnClick(func(e events.Event) {
						if deck.onEdit != nil {
							deck.onEdit()
						}
					})
			})
			i.OnClick(func(e events.Event) {
				i.ShowContextMenu(e)
			})
		})
	})
}
func (deck *Deck) IconButton(ctx *core.Frame, icon icons.Icon, callback func(i *core.Icon)) {
	tree.AddChildAt(ctx, strconv.Itoa(deck.index)+"-icon", func(w *core.Icon) {
		w.SetIcon(icon)
		w.Styler(func(s *styles.Style) {
			s.SetAbilities(true, abilities.Clickable)
			s.SetAbilities(true, abilities.Hoverable)
			s.Color = colors.Scheme.OnSurfaceVariant
			s.Font.Size.Set(20, units.UnitDp)
			if s.Is(states.Hovered) {
				s.Color = colors.Uniform(colors.Aliceblue)
				s.Color = colors.Scheme.OnSurface
			}
		})
		if callback != nil {
			callback(w)
		}

	})
	deck.index++
}
func (deck *Deck) OnExplore(f func()) {
	deck.onExplore = f
}

func (deck *Deck) OnAddCard(f func()) {
	deck.onAddCard = f
}

func (deck *Deck) OnEdit(f func()) {
	deck.onEdit = f
}

func (deck *Deck) OnStudy(f func()) {
	deck.onStudy = f
}

func (deck *Deck) OnMore(f func()) {
	deck.onMore = f
}

func (deck *Deck) OnDelete(f func()) {
	deck.onDelete = f
}
func (deck *Deck) setData(deckdata *models.Deck) {
	deck.deckdata = deckdata
}
