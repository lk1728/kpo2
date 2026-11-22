package application

import "pos/internal/domain"

type Service struct {
	repository CourseRepository
	exporter   ResultExporter
}

func NewService(repository CourseRepository, exporter ResultExporter) *Service {
	return &Service{
		repository: repository,
		exporter:   exporter,
	}
}

func (s *Service) AddStudent(name string) (domain.Student, error) {
	course, err := s.repository.Get()
	if err != nil {
		return domain.Student{}, err
	}

	student, err := course.AddStudent(name)
	if err != nil {
		return domain.Student{}, err
	}

	err = s.repository.Save(course)
	if err != nil {
		return domain.Student{}, err
	}

	return student, nil
}

func (s *Service) Students() ([]domain.Student, error) {
	course, err := s.repository.Get()
	if err != nil {
		return nil, err
	}

	return course.Students(), nil
}

func (s *Service) AddCriterion(name string, maxScore float64, weight float64) (domain.Criterion, error) {
	course, err := s.repository.Get()
	if err != nil {
		return domain.Criterion{}, err
	}

	criterion, err := course.AddCriterion(name, maxScore, weight)
	if err != nil {
		return domain.Criterion{}, err
	}

	err = s.repository.Save(course)
	if err != nil {
		return domain.Criterion{}, err
	}

	return criterion, nil
}

func (s *Service) Criteria() ([]domain.Criterion, error) {
	course, err := s.repository.Get()
	if err != nil {
		return nil, err
	}

	return course.Criteria(), nil
}

func (s *Service) SetEvaluation(studentID int, criterionID int, score float64) error {
	course, err := s.repository.Get()
	if err != nil {
		return err
	}

	err = course.SetEvaluation(studentID, criterionID, score)
	if err != nil {
		return err
	}

	return s.repository.Save(course)
}

func (s *Service) Evaluations(studentID int) ([]domain.Evaluation, error) {
	course, err := s.repository.Get()
	if err != nil {
		return nil, err
	}

	return course.Evaluations(studentID)
}

func (s *Service) RemoveEvaluation(studentID int, criterionID int) error {
	course, err := s.repository.Get()
	if err != nil {
		return err
	}

	err = course.RemoveEvaluation(studentID, criterionID)
	if err != nil {
		return err
	}

	return s.repository.Save(course)
}

func (s *Service) CalculateTotal(studentID int, current bool) (StudentResult, error) {
	course, err := s.repository.Get()
	if err != nil {
		return StudentResult{}, err
	}

	student, err := course.Student(studentID)
	if err != nil {
		return StudentResult{}, err
	}

	var strategy domain.ScoreStrategy = domain.FinalScoreStrategy{}
	if current {
		strategy = domain.CurrentScoreStrategy{}
	}

	calculator := domain.NewTotalCalculator(strategy)

	score, err := calculator.Calculate(course, studentID)
	if err != nil {
		return StudentResult{}, err
	}

	return StudentResult{
		StudentID:   student.ID(),
		StudentName: student.Name(),
		Score:       score,
	}, nil
}

func (s *Service) Results() ([]StudentResult, error) {
	course, err := s.repository.Get()
	if err != nil {
		return nil, err
	}

	calculator := domain.NewTotalCalculator(domain.FinalScoreStrategy{})
	results := make([]StudentResult, 0)

	for _, student := range course.Students() {
		score, err := calculator.Calculate(course, student.ID())
		if err != nil {
			return nil, err
		}

		results = append(results, StudentResult{
			StudentID:   student.ID(),
			StudentName: student.Name(),
			Score:       score,
		})
	}

	return results, nil
}

func (s *Service) ExportResults(path string) error {
	results, err := s.Results()
	if err != nil {
		return err
	}

	return s.exporter.Export(path, results)
}
