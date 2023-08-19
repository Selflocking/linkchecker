package utils

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v53/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"linkchecker/config"
)

var GitHubClient *github.Client

func init() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(tc.Transport)
	if err != nil {
		panic(err)
	}
	GitHubClient = github.NewClient(rateLimiter)
}

func GetAllRepos(org string) (repos []*github.Repository) {
	nowPage := 1
	for {
		rs, resp, err := GitHubClient.Repositories.List(
			context.Background(),
			org,
			&github.RepositoryListOptions{
				Visibility: "public",
				Type:       "public",
				ListOptions: github.ListOptions{
					Page: nowPage,
				},
			})
		if err != nil {
			log.Fatal(err)
		}

		for _, r := range rs {
			ok := true
			for _, ignoreRepo := range config.IgnoreRepos {
				if fmt.Sprintf("%s/%s", *r.Owner.Login, *r.Name) == ignoreRepo {
					ok = false
					break
				}
			}
			if ok && !*r.Fork && !*r.Archived {
				repos = append(repos, r)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		nowPage = nowPage + 1
	}

	return repos
}

func ForkRepo(owner string, repoName string, author string) {
	_, resp, _ := GitHubClient.Repositories.Get(context.Background(), author, repoName)
	if resp.StatusCode == 200 {
		return
	}

	_, _, _ = GitHubClient.Repositories.CreateFork(context.Background(), owner,
		repoName, &github.RepositoryCreateForkOptions{})
}

func OpenPR(owner string, repoName string, user string, branch string, prTitle string,
	prContent string) {
	_, _, err := GitHubClient.PullRequests.Create(context.Background(), owner,
		repoName,
		&github.NewPullRequest{
			Title:               github.String(prTitle),
			Head:                github.String(user + ":" + branch),
			Base:                github.String("master"),
			Body:                github.String(prContent),
			MaintainerCanModify: github.Bool(true),
		})
	if err != nil {
		logrus.Error("create PR failed: %s/%s\n,%v\n", owner, repoName, err)
		return
	}
}

func ParseName(repo string) (owner string, repoName string) {
	s := strings.Split(repo, "/")
	return s[0], s[1]
}

func IsUpdated(repo *github.Repository) bool {
	// check if updated in a week
	if repo.UpdatedAt.After(time.Now().AddDate(0, 0, -7)) {
		return true
	}
	return false
}
