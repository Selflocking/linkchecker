package utils

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func Clone(repoName string, p string, cloneUrl string) {
	if _, err := os.Stat(filepath.Join(p, repoName)); err == nil {
		return
		// cmd := exec.Command("git", "pull")
		// cmd.Dir = filepath.Join(p, repoName)
		// err = cmd.Run()
		// return err
	}

	if _, err := os.Stat(p); os.IsNotExist(err) {
		err = os.MkdirAll(p, os.ModePerm)
		if err != nil {
			logrus.Error(err)
			return
		}
	}

	cmd := exec.Command("git", "clone", cloneUrl, "--depth=1")
	cmd.Dir = p
	err := cmd.Run()

	if err != nil {
		logrus.Error("clone failed:", repoName)
		return
	}

	return
}

func AddRemote(p string, author string, repoName string) {
	getUrl := exec.Command("git", "remote", "get-url", "me")
	getUrl.Dir = p
	url, err := getUrl.Output()
	if err != nil {
		cmd := exec.Command("git", "remote", "add", "me", "git@github.com:"+author+"/"+repoName+".git")
		cmd.Dir = p
		output, err := cmd.Output()
		if err != nil {
			logrus.Error("In "+p+":", err)
			logrus.Error(string(output))
			return
		}
	}

	if string(url) != "git@github.com:"+author+"/"+repoName+".git" {
		cmd := exec.Command("git", "remote", "set-url", "me",
			"git@github.com:"+author+"/"+repoName+".git")
		cmd.Dir = p
		output, err := cmd.Output()
		if err != nil {
			logrus.Error("In "+p+":", err)
			logrus.Error(string(output))
			return
		}
	}
}

func CommitAndPush(p string, branch string, msg string, remote string) bool {
	checkout := exec.Command("git", "checkout", "-b", branch)
	checkout.Dir = p
	output, err := checkout.Output()
	if err != nil {
		logrus.Error("In "+p+":", err)
		logrus.Error(string(output))
		return false
	}

	add := exec.Command("git", "add", ".")
	add.Dir = p
	output, err = add.Output()
	if err != nil {
		logrus.Error("In "+p+":", err)
		logrus.Error(string(output))
		return false
	}

	commit := exec.Command("git", "commit", "-m", msg)
	commit.Dir = p
	output, err = commit.Output()
	if err != nil {
		logrus.Error("In "+p+":", err)
		logrus.Error(string(output))
		return false
	}

	push := exec.Command("git", "push", remote)
	push.Dir = p
	output, err = push.Output()
	if err != nil {
		logrus.Error("In "+p+":", err)
		logrus.Error(string(output))
		return false
	}

	back := exec.Command("git", "checkout", "master")
	back.Dir = p
	output, err = back.Output()
	if err != nil {
		logrus.Error("In "+p+":", err)
		logrus.Error(string(output))
		return false
	}
	return true
}

func OpenRepo(p string) {
	cmd := exec.Command("code", p)
	err := cmd.Run()
	if err != nil {
		logrus.Error(err)
		return
	}
}
