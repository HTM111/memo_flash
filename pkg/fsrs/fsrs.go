package fsrs

import (
	"log"
	"math"
	"memoflash/internal/models"
	"memoflash/internal/values"
	"time"
)

const (
	// Initial stability for new cards
	W0 = 0.4 // Again
	W1 = 0.6 // Hard
	W2 = 2.4 // Good
	W3 = 5.8 // Easy

	// Stability increase factors
	W4  = 4.93 // Good review multiplier
	W5  = 0.94 // Easy bonus multiplier
	W6  = 0.86 // Hard penalty multiplier
	W7  = 0.01 // Difficulty decay on success
	W8  = 1.49 // Difficulty increase on failure
	W9  = 0.14 // Retrievability factor
	W10 = 0.94 // Stability recovery factor

	MinStability  = 0.1
	MaxStability  = 36500 // ~100 years in days
	MinDifficulty = 1.0
	MaxDifficulty = 10.0

	RequestedRetention = 0.9
)

type FSRSEngine struct {
	card *models.Card
}

func Review(difficulty values.Difficulty, card *models.Card) {
	daysSinceLastReview := time.Since(card.LastStudied).Hours() / 24

	// Minimum time between reviews
	if daysSinceLastReview < 0.05 {
		daysSinceLastReview = 0.05
	}

	isNewCard := card.LastStudied.IsZero()

	switch difficulty {
	case values.Again:
		handleAgainReview(card, daysSinceLastReview, isNewCard)
	case values.Hard:
		handleHardReview(card, daysSinceLastReview, isNewCard)
	case values.Good:
		handleGoodReview(card, daysSinceLastReview, isNewCard)
	case values.Easy:
		handleEasyReview(card, daysSinceLastReview, isNewCard)
	default:
		log.Println("Invalid difficulty")
		return
	}

	// Apply bounds
	card.Stability = math.Max(MinStability, math.Min(MaxStability, card.Stability))
	card.Difficulty = math.Max(MinDifficulty, math.Min(MaxDifficulty, card.Difficulty))

	// Calculate next review interval
	interval := calculateInterval(card.Stability)
	card.Interval = time.Now().AddDate(0, 0, int(interval))
	card.LastStudied = time.Now()
}

func handleAgainReview(card *models.Card, daysSince float64, isNewCard bool) {
	if isNewCard {
		card.Stability = W0
		card.Difficulty = MinDifficulty
	} else {
		// For failed reviews, use actual retrievability to adjust stability decrease
		retrievability := calculateRetrievability(daysSince, card.Stability)
		// Worse retrievability (forgot sooner) = bigger stability decrease
		stabilityDecrease := W10 * (1.0 + (1.0 - retrievability)) * math.Pow(card.Difficulty/MaxDifficulty, W9)
		card.Stability = card.Stability * stabilityDecrease
		card.Difficulty = math.Min(card.Difficulty+W8, MaxDifficulty)
	}
}

func handleHardReview(card *models.Card, daysSince float64, isNewCard bool) {
	if isNewCard {
		card.Stability = W1
		card.Difficulty = MinDifficulty + 1.0
	} else {
		retrievability := calculateRetrievability(daysSince, card.Stability)
		stabilityIncrease := calculateStabilityIncrease(retrievability, card.Difficulty)

		card.Stability = card.Stability * stabilityIncrease * W6          // Hard penalty
		card.Difficulty = math.Max(card.Difficulty-W7*0.5, MinDifficulty) // Small difficulty decrease
	}
}

func handleGoodReview(card *models.Card, daysSince float64, isNewCard bool) {
	if isNewCard {
		card.Stability = W2
		card.Difficulty = MinDifficulty + 0.5
	} else {
		retrievability := calculateRetrievability(daysSince, card.Stability)
		stabilityIncrease := calculateStabilityIncrease(retrievability, card.Difficulty)

		card.Stability = card.Stability * stabilityIncrease * W4
		card.Difficulty = math.Max(card.Difficulty-W7, MinDifficulty)
	}
}

func handleEasyReview(card *models.Card, daysSince float64, isNewCard bool) {
	if isNewCard {
		card.Stability = W3
		card.Difficulty = MinDifficulty
	} else {
		retrievability := calculateRetrievability(daysSince, card.Stability)
		stabilityIncrease := calculateStabilityIncrease(retrievability, card.Difficulty)

		card.Stability = card.Stability * stabilityIncrease * W4 * W5     // Easy bonus
		card.Difficulty = math.Max(card.Difficulty-W7*1.5, MinDifficulty) // Larger difficulty decrease
	}
}

func calculateRetrievability(daysSince, stability float64) float64 {
	return math.Exp(math.Log(RequestedRetention) * daysSince / stability)
}

func calculateStabilityIncrease(retrievability, difficulty float64) float64 {
	difficultyFactor := (MaxDifficulty - difficulty) / MaxDifficulty
	retrievabilityFactor := 1.0 + W9*(1.0-retrievability)

	return 1.0 + difficultyFactor*retrievabilityFactor
}

func calculateInterval(stability float64) float64 {
	interval := stability * math.Log(RequestedRetention) / math.Log(RequestedRetention)
	interval = stability // At stability days, retrievability ≈ 0.9

	return math.Max(1.0, interval)
}
