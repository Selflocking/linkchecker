package utils

import (
	"bufio"
	"fmt"
	"linkchecker/config"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type Link struct {
	Url     string
	Org     string
	Repo    string
	RelPath string
	Line    int
}

func (l Link) GetGithubLink() string {
	return fmt.Sprintf("https://github.com/%s/%s/blob/master/%s#L%d", l.Org, l.Repo, l.RelPath, l.Line)
}

func (l Link) GetFIlePath() string {
	return filepath.Join(l.Org, l.Repo, l.RelPath)
}

const linkPattern = "https?://(?:[a-zA-Z0-9\\-]+\\.?)+[^\\s>)\\]\"'`$]*"

func ParseLink(text string) []string {
	re := regexp.MustCompile(linkPattern)

	return re.FindAllString(text, -1)
}

func ExtractLinksFromFile(f File) []Link {
	file, err := os.Open(f.FilePath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	//buf := make([]byte, 0, 64*1024)
	//scanner.Buffer(buf, 1024*1024)

	var res []Link
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

			res = append(res, Link{
				Url:     l,
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

	return res
}
