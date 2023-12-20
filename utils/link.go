package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/work4dev/linkchecker/config"
)

type Location struct {
	Org     string
	Repo    string
	RelPath string
	Line    int
}

func (l Location) GetGithubLink() string {
	return fmt.Sprintf("https://github.com/%s/%s/blob/master/%s#L%d", l.Org, l.Repo, l.RelPath, l.Line)
}

func (l Location) GetFIlePath() string {
	return filepath.Join(l.Org, l.Repo, l.RelPath)
}

var LinksToCheck = make(map[string][]Location)

// Characters that cannot appear in links: blank , < > ) ] \ " ' ` $
const linkPattern = `https?://(?:[a-zA-Z0-9\-]+\.?)+[^\s,<>)\]\\"'` + "`" + `$]*`

func ParseLink(text string) []string {
	re := regexp.MustCompile(linkPattern)

	return re.FindAllString(text, -1)
}

func ExtractLinksFromFile(f File) {
	file, err := os.Open(f.FilePath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// buf := make([]byte, 0, 64*1024)
	// scanner.Buffer(buf, 1024*1024)

	lineCount := 0
	for scanner.Scan() {
		lineCount++
		links := ParseLink(scanner.Text())
		for _, l := range links {
			skip := false
			for _, pattern := range config.IgnoreURLPatterns {
				r := regexp.MustCompile(pattern)
				if r.MatchString(l) {
					skip = true
					break
				}
			}

			if skip {
				continue
			}

			LinksToCheck[l] = append(LinksToCheck[l], Location{
				Org:     f.Org,
				Repo:    f.Repo,
				RelPath: f.RelPath,
				Line:    lineCount,
			})

		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("failed to scan file", f.FilePath)
		log.Printf("error: %v", err)
	}

}
