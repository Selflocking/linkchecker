package utils

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/go-github/v53/github"
	"github.com/work4dev/linkchecker/config"
)

type File struct {
	FilePath string
	Org      string
	Repo     string
	RelPath  string
}

func GetFilesList(repo *github.Repository) (files []File) {
	_ = filepath.Walk(filepath.Join(config.Workspace, *repo.Owner.Login, *repo.Name), func(path string,
		info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		ok := false
		for _, fileType := range config.FileTypes {
			if strings.HasSuffix(path, fileType) {
				ok = true
				break
			}
		}

		for _, pattern := range config.IgnoreFilePatterns {
			r := regexp.MustCompile(pattern)
			if r.MatchString(info.Name()) {
				ok = false
				break
			}
		}

		for _, pattern := range config.IgnoreDirPatterns {
			r := regexp.MustCompile(pattern)
			if r.MatchString(filepath.Dir(path)) {
				ok = false
				break
			}
		}

		if ok {
			rel, err := filepath.Rel(filepath.Join(config.Workspace, *repo.Owner.Login, *repo.Name), path)
			if err != nil {
				log.Fatal(err)
			}
			files = append(files, File{
				FilePath: path,
				Org:      *repo.Owner.Login,
				Repo:     *repo.Name,
				RelPath:  rel,
			})
		}

		return nil
	})
	return files
}

func GetFileContent(f File) string {
	content, err := os.ReadFile(f.FilePath)
	if err != nil {
		log.Println("failed to read file", f.FilePath)
		log.Println(err)
		return ""
	}
	return string(content)
}

func GetFileExt(f File) string {
	return filepath.Ext(f.FilePath)
}
