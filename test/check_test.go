package test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/work4dev/linkchecker/utils"
)

func TestCheckLink(t *testing.T) {
	timeout := time.Second * 10

	invalidURL := "https://a.invaildURL.com"
	validURL := "https://github.com"

	log.Println("check github.com")
	ok, msg := utils.CheckLink(validURL, timeout)
	assert.EqualValues(t, true, ok)
	log.Println(msg)

	log.Println("check a.invaildURL.com")
	ok, msg = utils.CheckLink(invalidURL, timeout)
	assert.EqualValues(t, false, ok)
	log.Println(msg)
}
