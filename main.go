package main

import (
	"fmt"
	"github.com/google/go-github/v53/github"
	"github.com/schollz/progressbar/v3"
	"linkchecker/config"
	"linkchecker/utils"
	"time"
)

func main() {
	// 1. range orgs, get all repos
	var repos []*github.Repository
	for _, org := range config.Orgs {
		repos = append(repos, utils.GetAllRepos(org)...)
	}
	fmt.Println("repos: ", len(repos))

	// 2. clone all repos
	for _, repo := range repos {
		_ = utils.Clone(repo)
	}
	fmt.Println("clone done")

	// 3. get all allowed type files
	var files []utils.File
	for _, repo := range repos {
		files = append(files, utils.GetFilesList(repo)...)
	}
	fmt.Println("files: ", len(files))

	// 4. extract all links from files line by line
	for _, file := range files {
		utils.ExtractLinksFromFile(file)
	}
	fmt.Println("links: ", len(utils.LinksToCheck))

	bar := progressbar.NewOptions(len(utils.LinksToCheck),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("checking..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	// 5. check all links
	ch := make(chan int, 12)
	for url, loc := range utils.LinksToCheck {
		ch <- 1
		go func(u string, l []utils.Location) {
			ok, msg := utils.CheckLink(u, time.Second*10)
			if !ok {
				utils.AddToReport(u, l, msg)
			}
			_ = bar.Add(1)
			<-ch
		}(url, loc)
	}
}
