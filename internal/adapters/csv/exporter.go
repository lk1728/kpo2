package csv

import (
	"encoding/csv"
	"os"
	"strconv"

	"pos/internal/application"
)

type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Export(path string, results []application.StudentResult) error {
	records := [][]string{
		{"student_id", "student_name", "final_score"},
	}

	for _, result := range results {
		records = append(records, []string{
			strconv.Itoa(result.StudentID),
			result.StudentName,
			strconv.FormatFloat(result.Score, 'f', 2, 64),
		})
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	err = writer.WriteAll(records)
	closeErr := file.Close()

	if err != nil {
		return err
	}

	return closeErr
}
