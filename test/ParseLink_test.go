package test

import (
	"testing"

	"linkchecker/utils"
)

func TestParseLinks(t *testing.T) {
	file := utils.File{
		FilePath: "links.txt",
		Org:      "test",
		Repo:     "test",
		RelPath:  "links.txt",
	}
	utils.ExtractLinksFromFile(file)
	for url, loc := range utils.LinksToCheck {
		t.Logf("result: %s", url)
		if url != "https://abc.com" {
			t.Fatal("failed to parse line: ", loc[0].Line)
		}
	}
}
