package utils

import (
	"fmt"
	"os"
	"time"
)

var Report = []string{"GIthubLink, URL, Message"}

// AddToReport add a bad link to report
func AddToReport(l Link, msg string) {
	Report = append(Report, fmt.Sprintf("'%s', '%s', '%s'", l.GetGithubLink(), l.Url, msg))
}

// GenerateReport write report to csv file
func GenerateReport() {
	filename := "report-" + time.Now().Format("2006-01-02-15-04-05") + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for _, line := range Report {
		_, err := fmt.Fprintln(file, line)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
