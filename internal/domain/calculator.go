package domain

import "errors"

type ScoreStrategy interface {
	Calculate(criteria []Criterion, evaluations []Evaluation) (float64, error)
}

type FinalScoreStrategy struct{}

func (s FinalScoreStrategy) Calculate(criteria []Criterion, evaluations []Evaluation) (float64, error) {
	if len(criteria) == 0 {
		return 0, errors.New("в рубрикаторе нет критериев")
	}

	scores := make(map[int]float64)
	for _, evaluation := range evaluations {
		scores[evaluation.CriterionID()] = evaluation.Score().Value()
	}

	var total float64
	var totalWeight float64

	for _, criterion := range criteria {
		score := scores[criterion.ID()]
		total += score / criterion.MaxScore() * criterion.Weight()
		totalWeight += criterion.Weight()
	}

	return total / totalWeight * 10, nil
}

type CurrentScoreStrategy struct{}

func (s CurrentScoreStrategy) Calculate(criteria []Criterion, evaluations []Evaluation) (float64, error) {
	if len(criteria) == 0 {
		return 0, errors.New("в рубрикаторе нет критериев")
	}

	scores := make(map[int]float64)
	for _, evaluation := range evaluations {
		scores[evaluation.CriterionID()] = evaluation.Score().Value()
	}

	var total float64
	var totalWeight float64

	for _, criterion := range criteria {
		score, exists := scores[criterion.ID()]
		if !exists {
			continue
		}

		total += score / criterion.MaxScore() * criterion.Weight()
		totalWeight += criterion.Weight()
	}

	if totalWeight == 0 {
		return 0, errors.New("у студента пока нет оценок")
	}

	return total / totalWeight * 10, nil
}

type TotalCalculator struct {
	strategy ScoreStrategy
}

func NewTotalCalculator(strategy ScoreStrategy) *TotalCalculator {
	return &TotalCalculator{strategy: strategy}
}

func (c *TotalCalculator) Calculate(course *Course, studentID int) (float64, error) {
	if course == nil {
		return 0, errors.New("курс не задан")
	}
	if c.strategy == nil {
		return 0, errors.New("стратегия расчёта не задана")
	}

	evaluations, err := course.Evaluations(studentID)
	if err != nil {
		return 0, err
	}

	return c.strategy.Calculate(course.Criteria(), evaluations)
}
