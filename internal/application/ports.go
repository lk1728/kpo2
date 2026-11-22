package application

import "pos/internal/domain"

type CourseRepository interface {
	Get() (*domain.Course, error)
	Save(course *domain.Course) error
}

type StudentResult struct {
	StudentID   int
	StudentName string
	Score       float64
}

type ResultExporter interface {
	Export(path string, results []StudentResult) error
}
