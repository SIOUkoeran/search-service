package csvreader

import (
	"encoding/csv"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	"os"
)

// read csv method
func ReadCsv(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	defer file.Close()

	utf8Reader := transform.NewReader(file, korean.EUCKR.NewDecoder())

	reader := csv.NewReader(utf8Reader)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
