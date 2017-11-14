package gh

import (
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	homedir "github.com/mitchellh/go-homedir"
	"path"
	"os"
	"strings"
)

var client *github.Client

// Construct client lazily such that we don't complain about
// credentials unless needed.
func Client() *github.Client {
	if client == nil {
		client = newGithubClient()
	}
	return client
}

// Create new client. 
func newGithubClient() *github.Client {
	// Create a BasicAuthTransport if the user has these env var
	// configured
	var basicAuth *github.BasicAuthTransport
	username := os.Getenv("BZL_GH_USERNAME")
	password := os.Getenv("BZL_GH_PASSWORD")
	if username != "" && password != "" {
		basicAuth = &github.BasicAuthTransport{
			Username: strings.TrimSpace(username),
			Password: strings.TrimSpace(password),
		}
	}

	// Create a cache/transport implementation
	var cacheTransport *httpcache.Transport
	homeDir, err := homedir.Dir()
	if err == nil {
		cacheDir := path.Join(homeDir, ".cache", "bzl", "gh")
		cache := diskcache.New(cacheDir)
		cacheTransport = httpcache.NewTransport(cache)
	} else {
		// Couldn't get homedir, just use mem cache instead
		// though of limited utility...
		cacheTransport = httpcache.NewMemoryCacheTransport()
	}

	// Create a Client 
	if basicAuth != nil {
		basicAuth.Transport = cacheTransport
		return github.NewClient(basicAuth.Client())
	} else {
		return github.NewClient(cacheTransport.Client())
	}
}
