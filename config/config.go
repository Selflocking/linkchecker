package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// https://oa.casbin.com/api/get-issue
var (
	Orgs               []string
	IgnoreRepos        []string
	FileTypes          []string
	IgnoreFilePatterns []string
	IgnoreDirPatterns  []string
	IgnoreURLPatterns  []string
)

var (
	GitHubToken = ""
	Author      = ""
	Workspace   = ""
	UA          = ""
)

func init() {
	viper.SetDefault("GitHubToken", "")
	viper.SetDefault("Author", "")
	viper.SetDefault("Workspace", "workspace")
	viper.SetDefault("UA", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36 Edg/115.0.1901.203")
	viper.SetDefault("Orgs", []string{})
	viper.SetDefault("IgnoreRepos", []string{})
	viper.SetDefault("FileTypes", []string{})
	viper.SetDefault("IgnoreFilePatterns", []string{})
	viper.SetDefault("IgnoreDirPatterns", []string{})
	viper.SetDefault("IgnoreURLPatterns", []string{})

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	_ = viper.ReadInConfig()

	// if viper.ReadInConfig() != nil {
	//	err := viper.SafeWriteConfig()
	//	if err != nil {
	//		panic(err)
	//	}
	// }

	viper.SetEnvPrefix("LINKCHECKER")
	viper.AutomaticEnv()

	GitHubToken = viper.GetString("GitHubToken")

	if GitHubToken == "" {
		panic("please provide github token")
	}

	Workspace = viper.GetString("Workspace")

	if Workspace == "" {
		panic("please provide a valid workspace path")
	}

	Author = viper.GetString("Author")

	if Author == "" {
		panic("please provide your github username")
	}
	UA = viper.GetString("UA")

	logrus.Info("Author: ", Author)
	logrus.Info("Workspace: ", Workspace)
	logrus.Info("UA: ", UA)

	Orgs = viper.GetStringSlice("Orgs")
	IgnoreRepos = viper.GetStringSlice("IgnoreRepos")
	FileTypes = viper.GetStringSlice("FileTypes")
	IgnoreFilePatterns = viper.GetStringSlice("IgnoreFilePatterns")
	IgnoreDirPatterns = viper.GetStringSlice("IgnoreDirPatterns")
	IgnoreURLPatterns = viper.GetStringSlice("IgnoreURLPatterns")
}
