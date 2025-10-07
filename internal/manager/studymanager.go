package manager

import (
	"errors"
	"fmt"
	"memoflash/internal/models"
	"memoflash/internal/services"
	"memoflash/internal/values"
	"memoflash/pkg/fsrs"
	"sort"
	"strings"
	"sync"
	"time"
)

type StudyManager struct {
	cardsService      services.CardService
	decksService      services.DeckService
	decks             []*models.Deck
	mu                sync.Mutex
	selectedDeck      *models.Deck
	selectedDeckCards []*models.Card

	cachedDueCount int
	dayStreak      int
	progress       int
}

func NewStudyManager(cardsService services.CardService, decksService services.DeckService) (*StudyManager, error) {
	if cardsService == nil || decksService == nil {
		return nil, errors.New("services cannot be nil")
	}

	sm := &StudyManager{
		decks:        make([]*models.Deck, 0),
		cardsService: cardsService,
		mu:           sync.Mutex{},
		decksService: decksService,
	}

	if err := sm.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize StudyManager: %w", err)
	}

	return sm, nil
}
func (sm *StudyManager) SearchCards(query string) []*models.Card {
	if sm.selectedDeckCards == nil {
		return nil
	}

	if query == "" {
		return sm.selectedDeckCards
	}
	var cards []*models.Card
	queryLower := strings.ToLower(query)
	for _, card := range sm.selectedDeckCards {
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
func (sm *StudyManager) refreshProgress() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	progress, err := sm.cardsService.GetProgress()
	if err != nil {
		return fmt.Errorf("failed to refresh progress: %w", err)
	}
	sm.progress = progress
	return nil
}
func (sm *StudyManager) refreshDayStreak() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	dayStreak, err := sm.cardsService.GetStreak()
	if err != nil {
		return fmt.Errorf("failed to refresh day streak: %w", err)
	}
	sm.dayStreak = dayStreak
	return nil
}
func (sm *StudyManager) initialize() error {
	if err := sm.loadDecks(); err != nil {
		return err
	}
	if err := sm.refreshDayStreak(); err != nil {
		return err
	}

	if err := sm.refreshDueCount(); err != nil {
		return err
	}
	if err := sm.refreshProgress(); err != nil {
		return err
	}
	return nil
}
func (sm *StudyManager) ClearSelectedCards() {
	sm.selectedDeckCards = []*models.Card{}
}
func (sm *StudyManager) loadDecks() error {
	decks, err := sm.decksService.GetDecks()
	if err != nil {
		return fmt.Errorf("failed to load decks: %w", err)
	}

	sm.decks = decks
	return nil
}

func (sm *StudyManager) refreshDueCount() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	count, err := sm.cardsService.CountDueCards()
	if err != nil {
		return fmt.Errorf("failed to refresh due count: %w", err)
	}
	sm.cachedDueCount = count
	return nil
}

func (sm *StudyManager) GetDecks() []*models.Deck {
	decksCopy := make([]*models.Deck, len(sm.decks))
	copy(decksCopy, sm.decks)
	return decksCopy
}

func (sm *StudyManager) GetDueCount() int {
	if sm == nil {
		return 0
	}
	return sm.cachedDueCount
}

func (sm *StudyManager) GetSelectedDeck() *models.Deck {
	return sm.selectedDeck
}

func (sm *StudyManager) GetSelectedCards() []*models.Card {
	return sm.selectedDeckCards
}

func (sm *StudyManager) GetDayStreak() int {
	return sm.dayStreak
}

func (sm *StudyManager) GetProgress() int {

	return sm.progress
}

func (sm *StudyManager) CreateDeck(title, desc string, categoryIndex int) error {
	id, err := sm.decksService.CreateDeck(title, desc, categoryIndex)
	if err != nil {
		return err
	}

	newDeck := &models.Deck{
		ID:            id,
		Title:         title,
		Description:   desc,
		CategoryIndex: categoryIndex,
	}
	sm.decks = append(sm.decks, newDeck)
	return nil
}
func (sm *StudyManager) SelectDeck(deck *models.Deck) error {

	sm.selectedDeck = deck
	cards, err := sm.cardsService.GetCardsFromDeck(deck.ID)
	if err != nil {
		return err
	}
	sm.selectedDeckCards = cards
	return nil

}
func (sm *StudyManager) GetAllDueCards() ([]*models.Card, error) {
	return sm.cardsService.GetAllDueCards()
}
func (sm *StudyManager) UpdateDeck(deck *models.Deck, title, desc string, categoryIndex int) error {
	if err := sm.decksService.EditDeck(deck.ID, title, desc, categoryIndex); err != nil {
		return err
	}

	deck.Title = title
	deck.Description = desc
	deck.CategoryIndex = categoryIndex

	return nil
}

func (sm *StudyManager) DeleteDeck(deck *models.Deck) error {
	if err := sm.decksService.DeleteDeck(deck.ID); err != nil {
		return err
	}

	for i, d := range sm.decks {
		if d.ID == deck.ID {
			sm.decks = append(sm.decks[:i], sm.decks[i+1:]...)
			break
		}
	}

	if sm.selectedDeck != nil && sm.selectedDeck.ID == deck.ID {
		sm.selectedDeck = nil
		sm.selectedDeckCards = nil
	}

	sm.refreshDueCount()
	sm.refreshProgress()
	return nil
}

func (sm *StudyManager) UpdateReadTime(deck *models.Deck) error {
	if deck == nil {
		return errors.New("deck cannot be nil")
	}

	if err := sm.decksService.UpdateReadTime(deck.ID); err != nil {
		return err
	}

	deck.LastStudied = time.Now()
	return nil
}

func (sm *StudyManager) CreateCard(front, back string, parentdeck *models.Deck) error {
	if err := sm.cardsService.CreateCard(front, back, parentdeck.ID); err != nil {
		return err
	}
	if err := sm.RefreshTotalCardsCount(parentdeck); err != nil {
		return err
	}
	err := sm.RefreshDueCardsInDeck(parentdeck)
	if err != nil {
		return err
	}
	if err := sm.refreshStates(); err != nil {
		return err
	}

	return nil
}

func (sm *StudyManager) EditCard(card *models.Card, front, back string) error {

	if card == nil {
		return errors.New("card cannot be nil")
	}

	if err := sm.cardsService.EditCard(card.ID, front, back); err != nil {
		return fmt.Errorf("failed to edit card: %s", err.Error())
	}

	card.Front = front
	card.Back = back

	return nil
}

func (sm *StudyManager) RefreshTotalCardsCount(deck *models.Deck) error {
	count, err := sm.cardsService.GetTotalCardsInDeck(deck.ID)
	if err != nil {
		return fmt.Errorf("failed to update total cards in deck: %s", err.Error())
	}
	deck.TotalCards = count
	return nil
}
func (sm *StudyManager) RefreshDueCardsInDeck(deck *models.Deck) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	due, err := sm.cardsService.CountDueCardsFromDeck(deck.ID)
	if err != nil {

		return err
	}
	deck.DueCards = due
	return nil
}
func (sm *StudyManager) DeleteCard(card *models.Card) error {
	if card == nil {
		return errors.New("card cannot be nil")
	}

	if err := sm.cardsService.DeleteCard(card.ID); err != nil {
		return err
	}
	if sm.selectedDeckCards != nil {
		for i, c := range sm.selectedDeckCards {
			if c.ID == card.ID {
				sm.selectedDeckCards = append(sm.selectedDeckCards[:i], sm.selectedDeckCards[i+1:]...)
				break
			}
		}
	}

	err := sm.RefreshTotalCardsCount(sm.selectedDeck)
	if err != nil {
		return err
	}
	err = sm.RefreshDueCardsInDeck(sm.selectedDeck)
	if err != nil {
		return err
	}

	err = sm.refreshStates()
	if err != nil {
		return err
	}
	return nil
}

func (sm *StudyManager) UpdateInterval(card *models.Card, difficulty values.Difficulty) error {
	if card == nil {
		return errors.New("card cannot be nil")
	}

	fsrs.Review(difficulty, card)

	if err := sm.decksService.UpdateInterval(card.ID, card.Interval, card.Difficulty, card.Stability); err != nil {
		return err
	}
	return sm.refreshStates()

}

func (sm *StudyManager) GetDueCardsFromDeck(deck *models.Deck) ([]*models.Card, error) {
	return sm.cardsService.GetDueCardsFromDeck(deck.ID)
}

func (sm *StudyManager) GetRecentlyStudiedDecks() []*models.Deck {
	var recentDecks = make([]*models.Deck, 0)
	for _, deck := range sm.decks {
		if !deck.LastStudied.IsZero() {
			recentDecks = append(recentDecks, deck)
		}
	}

	sort.Slice(recentDecks, func(i, j int) bool {
		return recentDecks[i].LastStudied.After(recentDecks[j].LastStudied)
	})

	if len(recentDecks) > 3 {
		recentDecks = recentDecks[:3]
	}

	return recentDecks
}

func (sm *StudyManager) refreshStates() error {

	err := sm.refreshDueCount()
	if err != nil {
		return err
	}

	err = sm.refreshProgress()
	if err != nil {
		return err
	}

	return nil
}
