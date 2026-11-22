package memory

import (
	"errors"

	"pos/internal/domain"
)

type Repository struct {
	course *domain.Course
}

func NewRepository(course *domain.Course) *Repository {
	return &Repository{course: course}
}

func (r *Repository) Get() (*domain.Course, error) {
	if r.course == nil {
		return nil, errors.New("курс не найден")
	}

	return r.course, nil
}

func (r *Repository) Save(course *domain.Course) error {
	if course == nil {
		return errors.New("нельзя сохранить пустой курс")
	}

	r.course = course
	return nil
}
