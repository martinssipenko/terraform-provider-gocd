// Copyright 2017 Drew J. Sonne. All rights reserved.
//
// Use of this source code is governed by a LGPLv3-style
// license that can be found in the LICENSE file.

/*
Package gocd provides a client for using the GoCD Server API.

Usage:

	import "github.com/drewsonne/go-gocd/gocd"

Construct a new GoCD client, then use the various services on the client to
access different parts of the GoCD Server API. For example:

	cfg := &gocd.Configuration{
		Server:       "https://goserver:8154/go/",
		Username:     os.GetEnv("GOCD_USERNAME"),
		Password:     os.GetEnv("GOCD_PASSWORD"),
		SkipSslCheck: true,
	}

	client := cfg.Client()

	// list all agents
	orgs, _, err := client.Agents.List(context.Background())

The services of a client divide the API into logical chunks and correspond to
the structure of the GoCD API documentation at
https://api.gocd.org/current/.

*/
package gocd
