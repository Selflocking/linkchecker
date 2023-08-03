package main

import (
	"fmt"
	"github.com/google/go-github/v53/github"
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
	var links []utils.Link
	for _, file := range files {
		links = append(links, utils.ExtractLinksFromFile(file)...)
	}
	fmt.Println("links: ", len(links))

	// 5. check all links
	ch := make(chan int, 12)
	for _, link := range links {
		ch <- 1
		go func(link utils.Link) {
			ok, msg := utils.CheckLink(link.Url, time.Second*10)
			if !ok {
				fmt.Printf("filepath: %s#%d, url: %s msg: %s\n",
					link.GetFIlePath(),
					link.Line,
					link.Url,
					msg)
				utils.AddToReport(link, msg)
			}
			<-ch
		}(link)
	}
}
