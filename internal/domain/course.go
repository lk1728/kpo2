package domain

import "errors"

type Course struct {
	id          int
	students    []Student
	criteria    []Criterion
	evaluations []Evaluation
}

func NewCourse(id int) (*Course, error) {
	if id <= 0 {
		return nil, errors.New("ID курса должен быть положительным")
	}

	return &Course{id: id}, nil
}

func (c *Course) ID() int {
	return c.id
}

func (c *Course) AddStudent(name string) (Student, error) {
	id := len(c.students) + 1

	student, err := NewStudent(id, name)
	if err != nil {
		return Student{}, err
	}

	c.students = append(c.students, student)
	return student, nil
}

func (c *Course) AddCriterion(name string, maxScore float64, weight float64) (Criterion, error) {
	id := len(c.criteria) + 1

	criterion, err := NewCriterion(id, name, maxScore, weight)
	if err != nil {
		return Criterion{}, err
	}

	c.criteria = append(c.criteria, criterion)
	return criterion, nil
}

func (c *Course) Student(id int) (Student, error) {
	for _, student := range c.students {
		if student.ID() == id {
			return student, nil
		}
	}

	return Student{}, errors.New("студент не найден")
}

func (c *Course) Criterion(id int) (Criterion, error) {
	for _, criterion := range c.criteria {
		if criterion.ID() == id {
			return criterion, nil
		}
	}

	return Criterion{}, errors.New("критерий не найден")
}

func (c *Course) Students() []Student {
	students := make([]Student, len(c.students))
	copy(students, c.students)
	return students
}

func (c *Course) Criteria() []Criterion {
	criteria := make([]Criterion, len(c.criteria))
	copy(criteria, c.criteria)
	return criteria
}

func (c *Course) SetEvaluation(studentID int, criterionID int, value float64) error {
	evaluation, err := NewEvaluation(c, studentID, criterionID, value)
	if err != nil {
		return err
	}

	for i, current := range c.evaluations {
		if current.StudentID() == studentID && current.CriterionID() == criterionID {
			c.evaluations[i] = evaluation
			return nil
		}
	}

	c.evaluations = append(c.evaluations, evaluation)
	return nil
}

func (c *Course) Evaluations(studentID int) ([]Evaluation, error) {
	_, err := c.Student(studentID)
	if err != nil {
		return nil, err
	}

	evaluations := make([]Evaluation, 0)

	for _, evaluation := range c.evaluations {
		if evaluation.StudentID() == studentID {
			evaluations = append(evaluations, evaluation)
		}
	}

	return evaluations, nil
}

func (c *Course) RemoveEvaluation(studentID int, criterionID int) error {
	_, err := c.Student(studentID)
	if err != nil {
		return err
	}

	_, err = c.Criterion(criterionID)
	if err != nil {
		return err
	}

	for i, evaluation := range c.evaluations {
		if evaluation.StudentID() == studentID && evaluation.CriterionID() == criterionID {
			c.evaluations = append(c.evaluations[:i], c.evaluations[i+1:]...)
			return nil
		}
	}

	return errors.New("оценка не найдена")
}
