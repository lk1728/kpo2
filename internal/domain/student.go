package domain

import (
	"errors"
	"strings"
)

type Student struct {
	id   int
	name string
}

func NewStudent(id int, name string) (Student, error) {
	name = strings.TrimSpace(name)
	if id <= 0 {
		return Student{}, errors.New("ID студента должен быть положительным")
	}
	if name == "" {
		return Student{}, errors.New("имя студента не должно быть пустым")
	}

	return Student{id: id, name: name}, nil
}

func (s Student) ID() int {
	return s.id
}

func (s Student) Name() string {
	return s.name
}
