package config

import (
	"path"
	"github.com/mitchellh/go-homedir"
)

var home = ""

func GetHome() (string, error) {
	if (home == "") {
		homeDir, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		home = path.Join(homeDir, ".cache", "bzl")
	}
	return home, nil
}
