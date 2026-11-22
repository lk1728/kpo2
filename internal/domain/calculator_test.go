package domain_test

import (
	"math"
	"testing"

	"pos/internal/domain"
)

func newCalculatorTestCourse(t *testing.T) *domain.Course {
	t.Helper()

	course, err := domain.NewCourse(1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = course.AddStudent("Николай")
	if err != nil {
		t.Fatal(err)
	}

	_, err = course.AddCriterion("Домашняя работа", 10, 2)
	if err != nil {
		t.Fatal(err)
	}

	_, err = course.AddCriterion("Экзамен", 20, 3)
	if err != nil {
		t.Fatal(err)
	}

	return course
}

func TestCalculatorStrategies(t *testing.T) {
	tests := []struct {
		name      string
		strategy  domain.ScoreStrategy
		scores    map[int]float64
		want      float64
		wantError bool
	}{
		{
			name:     "final with all grades",
			strategy: domain.FinalScoreStrategy{},
			scores:   map[int]float64{1: 8, 2: 20},
			want:     9.2,
		},
		{
			name:     "current with all grades",
			strategy: domain.CurrentScoreStrategy{},
			scores:   map[int]float64{1: 8, 2: 20},
			want:     9.2,
		},
		{
			name:     "final with missing grade",
			strategy: domain.FinalScoreStrategy{},
			scores:   map[int]float64{1: 8},
			want:     3.2,
		},
		{
			name:     "current with missing grade",
			strategy: domain.CurrentScoreStrategy{},
			scores:   map[int]float64{1: 8},
			want:     8,
		},
		{
			name:     "final without grades",
			strategy: domain.FinalScoreStrategy{},
			want:     0,
		},
		{
			name:      "current without grades",
			strategy:  domain.CurrentScoreStrategy{},
			wantError: true,
		},
		{
			name:     "current with zero grade",
			strategy: domain.CurrentScoreStrategy{},
			scores:   map[int]float64{1: 0},
			want:     0,
		},
		{
			name:     "current includes zero grade weight",
			strategy: domain.CurrentScoreStrategy{},
			scores:   map[int]float64{1: 0, 2: 20},
			want:     6,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			course := newCalculatorTestCourse(t)

			for criterionID, score := range test.scores {
				err := course.SetEvaluation(1, criterionID, score)
				if err != nil {
					t.Fatal(err)
				}
			}

			calculator := domain.NewTotalCalculator(test.strategy)
			got, err := calculator.Calculate(course, 1)

			if test.wantError {
				if err == nil {
					t.Fatal("ожидалась ошибка")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}
			if math.IsNaN(got) || math.IsInf(got, 0) {
				t.Fatalf("некорректный итог: %v", got)
			}
			if math.Abs(got-test.want) > 1e-9 {
				t.Fatalf("ожидалось %.2f, получено %.2f", test.want, got)
			}
		})
	}
}

func TestCalculatorRejectsInvalidInput(t *testing.T) {
	course := newCalculatorTestCourse(t)

	tests := []struct {
		name      string
		course    *domain.Course
		studentID int
		strategy  domain.ScoreStrategy
	}{
		{"missing course", nil, 1, domain.FinalScoreStrategy{}},
		{"unknown student", course, 99, domain.FinalScoreStrategy{}},
		{"missing strategy", course, 1, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calculator := domain.NewTotalCalculator(test.strategy)

			_, err := calculator.Calculate(test.course, test.studentID)
			if err == nil {
				t.Fatal("ожидалась ошибка")
			}
		})
	}
}

func TestCalculatorRejectsEmptyRubric(t *testing.T) {
	course, err := domain.NewCourse(1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = course.AddStudent("Николай")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		strategy domain.ScoreStrategy
	}{
		{"final", domain.FinalScoreStrategy{}},
		{"current", domain.CurrentScoreStrategy{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calculator := domain.NewTotalCalculator(test.strategy)

			_, err := calculator.Calculate(course, 1)
			if err == nil {
				t.Fatal("без критериев ожидалась ошибка")
			}
		})
	}
}
