package utils

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	ok  bool
	msg string
}

var CheckedLinks = make(map[string]Result)

// CheckLinkWithMemory Check if the link is available, and return directly if checked
func CheckLinkWithMemory(url string, timeout time.Duration) (ok bool, msg string) {
	if v, ok := CheckedLinks[url]; ok {
		return v.ok, v.msg
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ok = false
		msg = err.Error()
		CheckedLinks[url] = Result{ok, msg}
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ok = false
		msg = err.Error()
		CheckedLinks[url] = Result{ok, msg}
		return
	}
	if resp.StatusCode != 200 {
		ok = false
		msg = fmt.Sprintf("status code: %d", resp.StatusCode)
		CheckedLinks[url] = Result{ok, msg}
		return
	}

	ok = true
	msg = "ok"
	CheckedLinks[url] = Result{ok, msg}
	return ok, msg
}

// CheckLink Checks if the link is valid.
func CheckLink(url string, timeout time.Duration) (bool, string) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, err.Error()
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err.Error()
	}
	if resp.StatusCode != 200 {
		return false, fmt.Sprintf("status code: %d", resp.StatusCode)
	}

	return true, "ok"
}
