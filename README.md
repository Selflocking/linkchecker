# Link Checker

## Description

This is a program that checks the links in multiple repositories and generates a report.

## Usage

1. apply a GitHub token: https://github.com/settings/tokens
2. set an environment variable: `export LINKCHECKER_GITHUBTOKEN=your_token`
3. set an environment variable: `export LINKCHECKER_WORKSPACE=path/to/workspace`, workspace is a directory that clone all repositories you want to check.
4. run `go run main.go` or `go build && ./linkchecker`
5. If you need to fix broken links, set an environment variable: `export LINKCHECKER_AUTHOR=your_github_username`.
