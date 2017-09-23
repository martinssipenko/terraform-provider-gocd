package cli

import (
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// List of command name and descriptions
const (
	ConfigureCommandName  = "configure"
	ConfigureCommandUsage = "Generate configuration file ~/.gocd.conf"
)

func configureAction(c *cli.Context) (err error) {
	var cfg *gocd.Configuration
	var profile string

	if profile = c.Parent().String("profile"); profile == "" {
		profile = "default"
	}

	cfgs, err := gocd.LoadConfigFromFile()

	if cfg, err = generateConfig(); err != nil {
		return handleErrOutput("Configure:generate", err)
	} else {
		cfgs[profile] = cfg
	}

	b, err := yaml.Marshal(cfgs)
	if err != nil {
		return handleErrOutput("Configure:yaml", err)
	}

	if err = ioutil.WriteFile(gocd.ConfigFilePath(), b, 0644); err != nil {
		return handleErrOutput("Configure:write", err)
	}

	return nil
}

// Build a default template
func generateConfig() (cfg *gocd.Configuration, err error) {
	cfg = &gocd.Configuration{}
	qs := []*survey.Question{
		{
			Name:     "server",
			Prompt:   &survey.Input{Message: "GoCD Server (should contain '/go/' suffix)"},
			Validate: survey.Required,
		},
		{
			Name:   "username",
			Prompt: &survey.Input{Message: "Client Username"},
		},
		{
			Name:   "password",
			Prompt: &survey.Password{Message: "Client Password"},
		},
		{
			Name:   "skip_ssl_check",
			Prompt: &survey.Confirm{Message: "Skip SSL certificate validation"},
		},
	}

	err = survey.Ask(qs, cfg)
	return cfg, err
}

// ConfigureCommand handles the interaction between the cli flags and the action handler for configure
func configureCommand() *cli.Command {
	return &cli.Command{
		Name:   ConfigureCommandName,
		Usage:  ConfigureCommandUsage,
		Action: configureAction,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "profile"},
		},
	}
}
