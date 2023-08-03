package utils

import (
	"context"
	"fmt"
	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
	"linkchecker/config"
	"log"
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
