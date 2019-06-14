package exam

import (
	"encoding/json"
	"net/http"
	"time"
)

func getAllExams() ([]Exam, error) {
	url := "https://exam.kolhos.chichasov.es/"

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var exams []Exam

	err = json.NewDecoder(resp.Body).Decode(&exams)

	if err != nil {
		return nil, err
	}

	fillDates(exams)

	return exams, nil
}

func fillDates(exams []Exam) {
	for i := range exams {
		e := exams[i]

		e.DateTime, _ = time.Parse("2/1/2006", e.Date)
		e.DateTime.Add(23 * time.Hour) // It should show up the date of the exam.

		exams[i] = e
	}
}

func filterByTags(exams []Exam, tags ...[]string) []Exam {
	var filtered []Exam

	for i := range exams {
		for _, taggroup := range tags {
			ok := true
			for _, tag := range taggroup {
				if !contains(exams[i].Tags, tag) {
					ok = false
					break
				}
			}

			if ok {
				filtered = append(filtered, exams[i])
				break
			}
		}
	}

	return filtered
}

func contains(s []string, e string) bool {
	for i := range s {
		if e == s[i] {
			return true
		}
	}

	return false
}

func filterByDate(exams []Exam, from, to time.Time) []Exam {
	var filtered []Exam

	for i := range exams {
		if exams[i].DateTime.After(from) && exams[i].DateTime.Before(to) {
			filtered = append(filtered, exams[i])
		}
	}

	return filtered
}
