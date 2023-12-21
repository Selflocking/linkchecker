package test

import (
	"testing"

	"github.com/work4dev/linkchecker/utils"
)

var Urls = []string{
	"https://github.com/casbin",
	"https://github.com/casbin/casbin",
	"https://github.com/casdoor",
}

var Locs = [][]utils.Location{
	{
		{
			Org:     "casbin",
			Repo:    "casbin",
			RelPath: "README.md",
			Line:    0,
		},
	},
	{
		{

			Org:     "casbin",
			Repo:    "casbin",
			RelPath: "README.md",
			Line:    1,
		},
	},
	{
		{
			Org:     "casdoor",
			Repo:    "casdoor",
			RelPath: "README.md",
			Line:    0,
		},
	},
}

var Msg = []string{
	"ok",
	"ok",
	"ok",
}

func TestWriteReports(t *testing.T) {
	for i := 0; i < 3; i++ {
		utils.AddToReport(Urls[i], Locs[i], Msg[i])
	}
	utils.WriteReports("csv")
	utils.WriteReports("markdown")
}
