package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/go-github/v53/github"
	"github.com/sirupsen/logrus"
	"linkchecker/config"
	"linkchecker/utils"
)

func CheckAllRepos() {
	// 1. range orgs, get all repos
	var repos []*github.Repository
	for _, org := range config.Orgs {
		repos = append(repos, utils.GetAllRepos(org)...)
	}
	fmt.Println("repos: ", len(repos))

	// 2. clone all repos
	for _, repo := range repos {
		utils.Clone(*repo.Name, filepath.Join(config.Workspace, *repo.Owner.Login), *repo.CloneURL)
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

	linksNum := len(utils.LinksToCheck)
	fmt.Println("links: ", linksNum)
	counter := 0

	// 5. check all links
	ch := make(chan int, 12)
	defer close(ch)
	wg := sync.WaitGroup{}
	for url, loc := range utils.LinksToCheck {
		counter++
		logger := logrus.WithField("Rate", fmt.Sprintf("%d/%d", counter, linksNum))
		ch <- 1
		wg.Add(1)
		go func(u string, l []utils.Location) {
			defer wg.Done()

			ok, msg := utils.CheckLink(u, time.Second*20)
			if !ok {
				utils.AddToReport(u, l, msg)
			}
			logger.Info("Checking: ", u)
			<-ch
		}(url, loc)
	}

	wg.Wait()

	utils.WriteReports("markdown")
}

func ModifyAndOpenPR() {
	fmt.Print("Whether to start repairing dead links (y/n): ")
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Println(err)
		return
	}
	if input != "y" {
		return
	}

	const limit = 5
	counter := 0
	utils.Reports.Range(func(key, value interface{}) bool {
		repo, _ := key.(string)
		owner, repoName := utils.ParseName(repo)
		p := filepath.Join(config.Workspace, repo)

		utils.ForkRepo(owner, repoName, config.Author)
		utils.AddRemote(p, config.Author, repoName)
		if counter >= limit {
			for {
				fmt.Print("Whether to continue to open repos (y/n): ")
				var input string
				_, err := fmt.Scanln(&input)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if input == "y" {
					break
				}
			}
			counter = 0
		}
		utils.OpenRepo(p)
		counter++
		return true
	})

	for {
		fmt.Print("Whether it has been processed(y/n): ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println(err)
			return
		}
		if input == "y" {
			break
		}
	}

	msg := "fix: fix broken links"

	utils.Reports.Range(func(key, value interface{}) bool {
		repo, _ := key.(string)
		owner, repoName := utils.ParseName(repo)
		p := filepath.Join(config.Workspace, repo)

		ok := utils.CommitAndPush(p, "link", msg, "me")
		if ok {
			utils.OpenPR(owner, repoName, config.Author, "link", msg, "")
		}
		return true
	})
}

func CheckReposUpdatedWithinWeek() {
	// 1. range orgs, get all repos
	var repos []*github.Repository
	for _, org := range config.Orgs {
		repos = append(repos, utils.GetAllRepos(org)...)
	}

	var updatedRepos []*github.Repository
	for _, repo := range repos {
		if utils.IsUpdated(repo) {
			updatedRepos = append(updatedRepos, repo)
		}
	}
	fmt.Println("repos: ", len(updatedRepos))

	// 2. clone all repos
	for _, repo := range updatedRepos {
		utils.Clone(*repo.Name, filepath.Join(config.Workspace, *repo.Owner.Login), *repo.CloneURL)
	}
	fmt.Println("clone done")

	// 3. get all allowed type files
	var files []utils.File
	for _, repo := range updatedRepos {
		files = append(files, utils.GetFilesList(repo)...)
	}
	fmt.Println("files: ", len(files))

	// 4. extract all links from files line by line
	for _, file := range files {
		utils.ExtractLinksFromFile(file)
	}

	linksNum := len(utils.LinksToCheck)
	fmt.Println("links: ", linksNum)
	counter := 0

	// 5. check all links
	ch := make(chan int, 12)
	defer close(ch)
	wg := sync.WaitGroup{}
	for url, loc := range utils.LinksToCheck {
		counter++
		logger := logrus.WithField("Rate", fmt.Sprintf("%d/%d", counter, linksNum))
		ch <- 1
		wg.Add(1)
		go func(u string, l []utils.Location) {
			defer wg.Done()

			ok, msg := utils.CheckLink(u, time.Second*20)
			if !ok {
				utils.AddToReport(u, l, msg)
			}
			logger.Info("Checking: ", u)
			<-ch
		}(url, loc)
	}

	wg.Wait()

	utils.WriteReports("markdown")
}

func main() {
	var onlyCheckReoposWithinWeek = flag.Bool("week", false, "only check repos within a week")
	flag.Parse()

	if *onlyCheckReoposWithinWeek {
		fmt.Println("only check repos updated within a week")
		CheckReposUpdatedWithinWeek()
		return
	}
	fmt.Println("check all repos and fix broken links")
	CheckAllRepos()
	ModifyAndOpenPR()
}
