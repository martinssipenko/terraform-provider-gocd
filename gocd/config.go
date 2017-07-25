package gocdprovider

import "github.com/drewsonne/gocdsdk"

type ClientConfig struct {
	BaseURL string
	Auth    gocd.Auth
}
