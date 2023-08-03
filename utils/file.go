package utils

import (
	"github.com/google/go-github/v53/github"
	"linkchecker/config"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type File struct {
	FilePath string
	Org      string
	Repo     string
	RelPath  string
}

func Clone(repo *github.Repository) error {
	if _, err := os.Stat(filepath.Join(config.Workspace, *repo.Owner.Login, *repo.Name)); err == nil {
		return nil
		//cmd := exec.Command("git", "pull")
		//cmd.Dir = filepath.Join(config.Workspace, *repo.Owner.Login, *repo.Name)
		//err = cmd.Run()
		//return err
	}

	if _, err := os.Stat(filepath.Join(config.Workspace, *repo.Owner.Login)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join(config.Workspace, *repo.Owner.Login), os.ModePerm)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("git", "clone", *repo.CloneURL, "--depth=1")
	cmd.Dir = filepath.Join(config.Workspace, *repo.Owner.Login)
	err := cmd.Run()

	if err != nil {
		log.Fatalf("clone failed: %s/%s\n", *repo.Owner.Login, *repo.Name)
		return err
	}

	return nil
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
