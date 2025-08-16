package main

import (
	"encoding/csv"
	"os"
)

// ReadContactsEmailsPhones reads the CSV file passed as argument
// and returns two slices: emails and phone numbers
func ReadContactsEmailsPhones(filePath string) ([]string, []string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	var emails []string
	var phones []string

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}
		if len(row) >= 3 {
			emails = append(emails, row[1])
			phones = append(phones, row[2])
		}
	}

	return emails, phones, nil
}
