package domain

import (
	"errors"
	"math"
)

type Score struct {
	value float64
}

func NewScore(value float64, maxScore float64) (Score, error) {
	if maxScore <= 0 || math.IsNaN(maxScore) || math.IsInf(maxScore, 0) {
		return Score{}, errors.New("некорректный максимальный балл")
	}
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return Score{}, errors.New("балл должен быть конечным числом")
	}
	if value < 0 || value > maxScore {
		return Score{}, errors.New("балл выходит за пределы критерия")
	}

	return Score{value: value}, nil
}

func (s Score) Value() float64 {
	return s.value
}

type Evaluation struct {
	studentID   int
	criterionID int
	score       Score
}

func NewEvaluation(course *Course, studentID int, criterionID int, value float64) (Evaluation, error) {
	if course == nil {
		return Evaluation{}, errors.New("курс не задан")
	}

	_, err := course.Student(studentID)
	if err != nil {
		return Evaluation{}, err
	}

	criterion, err := course.Criterion(criterionID)
	if err != nil {
		return Evaluation{}, err
	}

	score, err := NewScore(value, criterion.MaxScore())
	if err != nil {
		return Evaluation{}, err
	}

	return Evaluation{
		studentID:   studentID,
		criterionID: criterionID,
		score:       score,
	}, nil
}

func (e Evaluation) StudentID() int {
	return e.studentID
}

func (e Evaluation) CriterionID() int {
	return e.criterionID
}

func (e Evaluation) Score() Score {
	return e.score
}
