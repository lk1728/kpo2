package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"pos/internal/application"
)

type Menu struct {
	service *application.Service
	scanner *bufio.Scanner
}

func NewMenu(service *application.Service) *Menu {
	return &Menu{
		service: service,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (m *Menu) Run() error {
	for {
		fmt.Print(`
1. Добавить студента
2. Показать студентов
3. Добавить критерий
4. Показать критерии
5. Выставить или исправить оценку
6. Показать оценки студента
7. Отменить оценку
8. Рассчитать балл студента
9. Показать итоговые баллы всех студентов
10. Экспортировать итоговые баллы в CSV
0. Выход

`)

		command, err := m.readString("Выберите действие: ")
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		switch command {
		case "1":
			err = m.addStudent()
		case "2":
			err = m.showStudents()
		case "3":
			err = m.addCriterion()
		case "4":
			err = m.showCriteria()
		case "5":
			err = m.setEvaluation()
		case "6":
			err = m.showEvaluations()
		case "7":
			err = m.removeEvaluation()
		case "8":
			err = m.calculateTotal()
		case "9":
			err = m.showResults()
		case "10":
			err = m.exportResults()
		case "0":
			return nil
		default:
			fmt.Println("Такого пункта меню нет")
		}

		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			fmt.Println("Ошибка:", err)
		}
	}
}

func (m *Menu) readString(prompt string) (string, error) {
	fmt.Print(prompt)

	if !m.scanner.Scan() {
		if err := m.scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}

	return strings.TrimSpace(m.scanner.Text()), nil
}

func (m *Menu) readInt(prompt string) (int, error) {
	text, err := m.readString(prompt)
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, errors.New("нужно ввести целое число")
	}

	return value, nil
}

func (m *Menu) readFloat(prompt string) (float64, error) {
	text, err := m.readString(prompt)
	if err != nil {
		return 0, err
	}

	text = strings.ReplaceAll(text, ",", ".")

	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, errors.New("нужно ввести число")
	}

	return value, nil
}

func (m *Menu) addStudent() error {
	name, err := m.readString("Имя студента: ")
	if err != nil {
		return err
	}

	student, err := m.service.AddStudent(name)
	if err != nil {
		return err
	}

	fmt.Printf("Добавлен студент: %d — %s\n", student.ID(), student.Name())
	return nil
}

func (m *Menu) showStudents() error {
	students, err := m.service.Students()
	if err != nil {
		return err
	}

	if len(students) == 0 {
		fmt.Println("Студентов пока нет")
		return nil
	}

	for _, student := range students {
		fmt.Printf("%d — %s\n", student.ID(), student.Name())
	}

	return nil
}

func (m *Menu) addCriterion() error {
	name, err := m.readString("Название критерия: ")
	if err != nil {
		return err
	}

	maxScore, err := m.readFloat("Максимальный балл: ")
	if err != nil {
		return err
	}

	weight, err := m.readFloat("Вес критерия: ")
	if err != nil {
		return err
	}

	criterion, err := m.service.AddCriterion(name, maxScore, weight)
	if err != nil {
		return err
	}

	fmt.Printf("Добавлен критерий: %d — %s\n", criterion.ID(), criterion.Name())
	return nil
}

func (m *Menu) showCriteria() error {
	criteria, err := m.service.Criteria()
	if err != nil {
		return err
	}

	if len(criteria) == 0 {
		fmt.Println("Критериев пока нет")
		return nil
	}

	for _, criterion := range criteria {
		fmt.Printf(
			"%d — %s, максимум: %.2f, вес: %.2f\n",
			criterion.ID(),
			criterion.Name(),
			criterion.MaxScore(),
			criterion.Weight(),
		)
	}

	return nil
}

func (m *Menu) setEvaluation() error {
	studentID, err := m.readInt("ID студента: ")
	if err != nil {
		return err
	}

	criterionID, err := m.readInt("ID критерия: ")
	if err != nil {
		return err
	}

	score, err := m.readFloat("Балл: ")
	if err != nil {
		return err
	}

	err = m.service.SetEvaluation(studentID, criterionID, score)
	if err != nil {
		return err
	}

	fmt.Println("Оценка сохранена")
	return nil
}

func (m *Menu) showEvaluations() error {
	studentID, err := m.readInt("ID студента: ")
	if err != nil {
		return err
	}

	evaluations, err := m.service.Evaluations(studentID)
	if err != nil {
		return err
	}

	if len(evaluations) == 0 {
		fmt.Println("У студента пока нет оценок")
		return nil
	}

	for _, evaluation := range evaluations {
		fmt.Printf(
			"Критерий %d — балл: %.2f\n",
			evaluation.CriterionID(),
			evaluation.Score().Value(),
		)
	}

	return nil
}

func (m *Menu) removeEvaluation() error {
	studentID, err := m.readInt("ID студента: ")
	if err != nil {
		return err
	}

	criterionID, err := m.readInt("ID критерия: ")
	if err != nil {
		return err
	}

	err = m.service.RemoveEvaluation(studentID, criterionID)
	if err != nil {
		return err
	}

	fmt.Println("Оценка отменена")
	return nil
}

func (m *Menu) calculateTotal() error {
	studentID, err := m.readInt("ID студента: ")
	if err != nil {
		return err
	}

	fmt.Println("1. Итоговый балл по всем критериям")
	fmt.Println("2. Текущий балл по выставленным оценкам")

	mode, err := m.readInt("Режим расчёта: ")
	if err != nil {
		return err
	}
	if mode != 1 && mode != 2 {
		return errors.New("нужно выбрать режим 1 или 2")
	}

	result, err := m.service.CalculateTotal(studentID, mode == 2)
	if err != nil {
		return err
	}

	if mode == 1 {
		fmt.Printf("%s — итоговый балл: %.2f из 10\n", result.StudentName, result.Score)
	} else {
		fmt.Printf("%s — текущий балл: %.2f из 10\n", result.StudentName, result.Score)
	}

	return nil
}

func (m *Menu) showResults() error {
	results, err := m.service.Results()
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("Студентов пока нет")
		return nil
	}

	for _, result := range results {
		fmt.Printf(
			"%d — %s: %.2f из 10\n",
			result.StudentID,
			result.StudentName,
			result.Score,
		)
	}

	return nil
}

func (m *Menu) exportResults() error {
	path, err := m.readString("Путь к CSV-файлу: ")
	if err != nil {
		return err
	}

	err = m.service.ExportResults(path)
	if err != nil {
		return err
	}

	fmt.Println("Результаты сохранены:", path)
	return nil
}
