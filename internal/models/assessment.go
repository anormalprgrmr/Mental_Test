package models

import (
	"fmt"
	"sort"
	"strings"
)

// Option is one selectable answer for a question.
type Option struct {
	Text   string
	Scores map[string]int
}

// Question belongs to a test and contributes to one or more dimensions.
type Question struct {
	Text    string
	Options []Option
}

// Dimension describes a score category in a test.
type Dimension struct {
	ID          string
	Name        string
	Description string
}

// TestDefinition is data-only, so adding new tests only requires registering one more definition.
type TestDefinition struct {
	ID          string
	Title       string
	Description string
	Dimensions  []Dimension
	Questions   []Question
}

// RankedScore is a dimension sorted by score.
type RankedScore struct {
	Dimension   Dimension
	Score       int
	MaxPossible int
}

// Result contains the calculated output shown to users/admins and persisted as JSON.
type Result struct {
	TestID   string        `json:"test_id"`
	Title    string        `json:"title"`
	Scores   []RankedScore `json:"scores"`
	Summary  string        `json:"summary"`
	Answers  []int         `json:"answers"`
	MaxScore int           `json:"max_score"`
}

func (t TestDefinition) Validate() error {
	if t.ID == "" || t.Title == "" {
		return fmt.Errorf("test id and title are required")
	}
	if len(t.Questions) == 0 {
		return fmt.Errorf("test %q has no questions", t.ID)
	}

	dimensions := make(map[string]struct{}, len(t.Dimensions))
	for _, dimension := range t.Dimensions {
		if dimension.ID == "" || dimension.Name == "" {
			return fmt.Errorf("test %q has invalid dimension", t.ID)
		}
		dimensions[dimension.ID] = struct{}{}
	}

	for questionIndex, question := range t.Questions {
		if question.Text == "" || len(question.Options) == 0 {
			return fmt.Errorf("test %q question %d is invalid", t.ID, questionIndex+1)
		}
		for optionIndex, option := range question.Options {
			if option.Text == "" {
				return fmt.Errorf("test %q question %d option %d has no text", t.ID, questionIndex+1, optionIndex+1)
			}
			for dimensionID := range option.Scores {
				if _, ok := dimensions[dimensionID]; !ok {
					return fmt.Errorf("test %q question %d option %d references unknown dimension %q", t.ID, questionIndex+1, optionIndex+1, dimensionID)
				}
			}
		}
	}
	return nil
}

func (t TestDefinition) Calculate(answers []int) (Result, error) {
	if len(answers) != len(t.Questions) {
		return Result{}, fmt.Errorf("expected %d answers, got %d", len(t.Questions), len(answers))
	}

	scores := make(map[string]int, len(t.Dimensions))
	maxPossible := make(map[string]int, len(t.Dimensions))
	for _, dimension := range t.Dimensions {
		scores[dimension.ID] = 0
		maxPossible[dimension.ID] = 0
	}

	for questionIndex, answerIndex := range answers {
		question := t.Questions[questionIndex]
		if answerIndex < 0 || answerIndex >= len(question.Options) {
			return Result{}, fmt.Errorf("invalid answer %d for question %d", answerIndex, questionIndex+1)
		}

		for dimensionID, value := range question.Options[answerIndex].Scores {
			scores[dimensionID] += value
		}

		bestForQuestion := make(map[string]int)
		for _, option := range question.Options {
			for dimensionID, value := range option.Scores {
				if value > bestForQuestion[dimensionID] {
					bestForQuestion[dimensionID] = value
				}
			}
		}
		for dimensionID, value := range bestForQuestion {
			maxPossible[dimensionID] += value
		}
	}

	ranked := make([]RankedScore, 0, len(t.Dimensions))
	for _, dimension := range t.Dimensions {
		ranked = append(ranked, RankedScore{
			Dimension:   dimension,
			Score:       scores[dimension.ID],
			MaxPossible: maxPossible[dimension.ID],
		})
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].Score == ranked[j].Score {
			return ranked[i].Dimension.Name < ranked[j].Dimension.Name
		}
		return ranked[i].Score > ranked[j].Score
	})

	return Result{
		TestID:   t.ID,
		Title:    t.Title,
		Scores:   ranked,
		Summary:  buildSummary(t.Title, ranked),
		Answers:  append([]int(nil), answers...),
		MaxScore: totalMax(ranked),
	}, nil
}

func buildSummary(title string, scores []RankedScore) string {
	if len(scores) == 0 {
		return fmt.Sprintf("نتیجه آزمون %s آماده است.", title)
	}

	limit := 3
	if len(scores) < limit {
		limit = len(scores)
	}

	parts := make([]string, 0, limit)
	for _, score := range scores[:limit] {
		parts = append(parts, fmt.Sprintf("%s (%d/%d)", score.Dimension.Name, score.Score, score.MaxPossible))
	}
	return fmt.Sprintf("قوی‌ترین حوزه‌های شما در %s: %s.", title, strings.Join(parts, "، "))
}

func totalMax(scores []RankedScore) int {
	total := 0
	for _, score := range scores {
		total += score.MaxPossible
	}
	return total
}
