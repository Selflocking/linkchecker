package config

import (
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
	Workspace   = "workspace"
)

func init() {
	viper.SetDefault("GitHubToken", "")
	viper.SetDefault("Workspace", "workspace")
	viper.SetDefault("Orgs", []string{})
	viper.SetDefault("IgnoreRepos", []string{})
	viper.SetDefault("FileTypes", []string{})
	viper.SetDefault("IgnoreFilePatterns", []string{})
	viper.SetDefault("IgnoreDirPatterns", []string{})
	viper.SetDefault("IgnoreURLPatterns", []string{})

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if viper.ReadInConfig() != nil {
		err := viper.SafeWriteConfig()
		if err != nil {
			panic(err)
		}
	}

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

	Orgs = viper.GetStringSlice("Orgs")
	IgnoreRepos = viper.GetStringSlice("IgnoreRepos")
	FileTypes = viper.GetStringSlice("FileTypes")
	IgnoreFilePatterns = viper.GetStringSlice("IgnoreFilePatterns")
	IgnoreDirPatterns = viper.GetStringSlice("IgnoreDirPatterns")
	IgnoreURLPatterns = viper.GetStringSlice("IgnoreURLPatterns")
}
