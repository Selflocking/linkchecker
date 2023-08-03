package utils

import (
	"fmt"
	"os"
	"time"
)

var logFile *os.File

func init() {
	if logFile == nil {
		var err error
		filename := "report-" + time.Now().Format("2006-01-02-15-04-05") + ".csv"
		logFile, err = os.Create(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = fmt.Fprintln(logFile, "GItHubLink, URL, Message")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// AddToReport add a bad link to report
func AddToReport(url string, loc []Location, msg string) {
	for _, l := range loc {
		_, err := fmt.Fprintln(logFile, fmt.Sprintf("%s, %s, \"%s\"", l.GetGithubLink(), url, msg))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
