package gocd

import "github.com/drewsonne/go-gocd/gocd"

type ClientConfig struct {
	BaseURL string
	Auth    gocd.Auth
}
