package utils

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type Result struct {
	Url     string
	RelPath string
	Line    int
	msg     string
}

var Reports sync.Map

// AddToReport add a bad link to report
func AddToReport(url string, loc []Location, msg string) {
	for _, l := range loc {
		results, ok := Reports.Load(l.Org + "/" + l.Repo)

		if ok {
			results = append(results.([]Result), Result{
				Url:     url,
				RelPath: l.RelPath,
				Line:    l.Line,
				msg:     msg,
			})
			Reports.Store(l.Org+"/"+l.Repo, results)
		} else {
			rs := make([]Result, 0)
			rs = append(rs, Result{
				Url:     url,
				RelPath: l.RelPath,
				Line:    l.Line,
				msg:     msg,
			})
			Reports.Store(l.Org+"/"+l.Repo, rs)
		}
	}
}

func WriteReports(t string) {
	switch t {
	case "csv":
		WriteReportsToCSV()
	case "markdown":
		WriteReportsToMD()
	default:
		fmt.Println("unknown report type")
	}
}

func WriteReportsToCSV() {
	filename := "report-" + time.Now().Format("2006-01-02-15-04-05") + ".csv"

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("failed to create file: ", err)
		return
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, "Repo, URL, Path, Message")
	if err != nil {
		fmt.Println("failed to write: ", err)
		return
	}

	Reports.Range(func(key, value interface{}) bool {
		results, _ := value.([]Result)
		repo, _ := key.(string)
		for _, r := range results {
			_, err := fmt.Fprintf(file, "%s, %s, %s#%d, %s\n", repo, r.Url, r.RelPath, r.Line, r.msg)
			if err != nil {
				fmt.Println("failed to write: ", err)
				return false
			}
		}
		_, err := fmt.Fprintf(file, "\n")
		if err != nil {
			fmt.Println("failed to write: ", err)
			return false
		}
		return true
	})
}

func WriteReportsToMD() {
	filename := "report-" + time.Now().Format("2006-01-02-15-04-05") + ".md"

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("failed to create file: ", err)
		return
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "# Report\n\n")
	if err != nil {
		fmt.Println("failed to write: ", err)
		return
	}

	Reports.Range(func(key, value interface{}) bool {
		results, _ := value.([]Result)
		repo, _ := key.(string)

		_, err := fmt.Fprintf(file, "## "+repo+"\n\n")
		if err != nil {
			fmt.Println("failed to write: ", err)
			return false
		}

		for _, r := range results {
			_, err := fmt.Fprintf(file,
				"URL: `%s`\nPath: %s#%d\nMessage: %s\n\n",
				r.Url, r.RelPath, r.Line, r.msg)
			if err != nil {
				fmt.Println("failed to write: ", err)
				return false
			}
		}
		return true
	})
}
