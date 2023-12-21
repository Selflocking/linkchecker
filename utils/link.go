package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/work4dev/linkchecker/comment"
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

func ignoreLink(link string) bool {
	for _, pattern := range config.IgnoreURLPatterns {
		r := regexp.MustCompile(pattern)
		if r.MatchString(link) {
			return true
		}
	}
	return false
}

func linkFilter(links []string) []string {
	var res []string
	for _, l := range links {
		if ignoreLink(l) {
			continue
		}
		res = append(res, l)
	}
	return res
}

func appendLinkToCheck(link []string, loc Location) {
	for _, l := range link {
		LinksToCheck[l] = append(LinksToCheck[l], loc)
	}
}

func ExtractLinksFromComments(f File) {
	ext := GetFileExt(f)
	parser := comment.ParserFactory(ext)
	if parser == nil {
		ExtractLinksFromFile(f)
		return
	}

	text := GetFileContent(f)
	comments := parser.GetComments(text)
	for _, c := range comments {
		links := ParseLink(c)
		links = linkFilter(links)
		appendLinkToCheck(links, Location{
			Org:     f.Org,
			Repo:    f.Repo,
			RelPath: f.RelPath,
			Line:    -1,
		})
	}
}

func ExtractLinksFromFile(f File) {
	file, err := os.Open(f.FilePath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lineCount := 0
	for scanner.Scan() {
		lineCount++
		links := ParseLink(scanner.Text())
		links = linkFilter(links)
		appendLinkToCheck(links, Location{
			Org:     f.Org,
			Repo:    f.Repo,
			RelPath: f.RelPath,
			Line:    lineCount,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Println("failed to scan file", f.FilePath)
		log.Printf("error: %v", err)
	}
}
