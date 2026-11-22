package domain

import (
	"errors"
	"math"
	"strings"
)

type Criterion struct {
	id       int
	name     string
	maxScore float64
	weight   float64
}

func NewCriterion(id int, name string, maxScore float64, weight float64) (Criterion, error) {
	name = strings.TrimSpace(name)
	if id <= 0 {
		return Criterion{}, errors.New("ID критерия должен быть положительным")
	}
	if name == "" {
		return Criterion{}, errors.New("название критерия не должно быть пустым")
	}
	if maxScore <= 0 || math.IsNaN(maxScore) || math.IsInf(maxScore, 0) {
		return Criterion{}, errors.New("максимальный балл должен быть положительным конечным числом")
	}
	if weight <= 0 || math.IsNaN(weight) || math.IsInf(weight, 0) {
		return Criterion{}, errors.New("вес должен быть положительным конечным числом")
	}

	return Criterion{
		id:       id,
		name:     name,
		maxScore: maxScore,
		weight:   weight,
	}, nil
}

func (c Criterion) ID() int {
	return c.id
}

func (c Criterion) Name() string {
	return c.name
}

func (c Criterion) MaxScore() float64 {
	return c.maxScore
}

func (c Criterion) Weight() float64 {
	return c.weight
}
