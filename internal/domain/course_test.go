package domain_test

import (
	"math"
	"testing"

	"pos/internal/domain"
)

func newTestCourse(t *testing.T) *domain.Course {
	t.Helper()

	course, err := domain.NewCourse(1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = course.AddStudent("Николай")
	if err != nil {
		t.Fatal(err)
	}

	_, err = course.AddCriterion("Домашняя работа", 10, 0.4)
	if err != nil {
		t.Fatal(err)
	}

	_, err = course.AddCriterion("Экзамен", 20, 0.6)
	if err != nil {
		t.Fatal(err)
	}

	return course
}

func TestSetEvaluationRejectsInvalidData(t *testing.T) {
	tests := []struct {
		name        string
		studentID   int
		criterionID int
		score       float64
	}{
		{"negative score", 1, 1, -1},
		{"score above maximum", 1, 1, 11},
		{"unknown student", 99, 1, 5},
		{"unknown criterion", 1, 99, 5},
		{"NaN score", 1, 1, math.NaN()},
		{"infinite score", 1, 1, math.Inf(1)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			course := newTestCourse(t)

			err := course.SetEvaluation(test.studentID, test.criterionID, test.score)
			if err == nil {
				t.Fatal("ожидалась ошибка")
			}

			evaluations, err := course.Evaluations(1)
			if err != nil {
				t.Fatal(err)
			}
			if len(evaluations) != 0 {
				t.Fatal("некорректная оценка не должна сохраняться")
			}
		})
	}
}

func TestSetEvaluationUpdatesScore(t *testing.T) {
	course := newTestCourse(t)

	err := course.SetEvaluation(1, 1, 8)
	if err != nil {
		t.Fatal(err)
	}

	err = course.SetEvaluation(1, 1, 9)
	if err != nil {
		t.Fatal(err)
	}

	evaluations, err := course.Evaluations(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(evaluations) != 1 {
		t.Fatalf("ожидалась одна оценка, получено %d", len(evaluations))
	}
	if evaluations[0].Score().Value() != 9 {
		t.Fatal("оценка не обновилась")
	}

	err = course.SetEvaluation(1, 1, 11)
	if err == nil {
		t.Fatal("ожидалась ошибка")
	}

	evaluations, err = course.Evaluations(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(evaluations) != 1 || evaluations[0].Score().Value() != 9 {
		t.Fatal("неудачное исправление изменило сохранённую оценку")
	}
}

func TestRemoveEvaluation(t *testing.T) {
	course := newTestCourse(t)

	err := course.SetEvaluation(1, 1, 8)
	if err != nil {
		t.Fatal(err)
	}

	err = course.SetEvaluation(1, 2, 15)
	if err != nil {
		t.Fatal(err)
	}

	err = course.RemoveEvaluation(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	evaluations, err := course.Evaluations(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(evaluations) != 1 || evaluations[0].CriterionID() != 2 {
		t.Fatal("после удаления должна остаться оценка за экзамен")
	}

	err = course.RemoveEvaluation(1, 1)
	if err == nil {
		t.Fatal("при повторном удалении ожидалась ошибка")
	}
}

func TestAddStudentRejectsEmptyName(t *testing.T) {
	course := newTestCourse(t)

	_, err := course.AddStudent("   ")
	if err == nil {
		t.Fatal("ожидалась ошибка для пустого имени")
	}
	if len(course.Students()) != 1 {
		t.Fatal("студент с пустым именем не должен добавляться")
	}
}

func TestAddCriterionRejectsInvalidData(t *testing.T) {
	tests := []struct {
		name     string
		maxScore float64
		weight   float64
	}{
		{"", 10, 1},
		{"Критерий", 0, 1},
		{"Критерий", -10, 1},
		{"Критерий", 10, 0},
		{"Критерий", 10, -1},
		{"Критерий", math.NaN(), 1},
		{"Критерий", 10, math.Inf(1)},
	}

	for _, test := range tests {
		course := newTestCourse(t)

		_, err := course.AddCriterion(test.name, test.maxScore, test.weight)
		if err == nil {
			t.Fatalf("ожидалась ошибка для критерия %+v", test)
		}
		if len(course.Criteria()) != 2 {
			t.Fatal("некорректный критерий не должен добавляться")
		}
	}
}
