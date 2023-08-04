package test

import (
	"github.com/stretchr/testify/assert"
	"linkchecker/utils"
	"log"
	"testing"
	"time"
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
