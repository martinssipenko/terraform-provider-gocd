package cli

import (
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

// List of command name and descriptions
const (
	ConfigureCommandName  = "configure"
	ConfigureCommandUsage = "Generate configuration file ~/.gocd.conf"
)

// ConfigDirectoryPath is the location where the authentication information is stored
const ConfigDirectoryPath = "~/.gocd.conf"

// Environment variables for configuration.
const (
	EnvVarServer   = "GOCD_SERVER"
	EnvVarUsername = "GOCD_USERNAME"
	EnvVarPassword = "GOCD_PASSWORD"
	EnvVarSkipSsl  = "GOCD_SKIP_SSL_CHECK"
)

//
func configureAction(c *cli.Context) error {

	s, err := generateConfigFile()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configFilePath(), []byte(s), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Build a default template
func generateConfigFile() (string, error) {
	cfg := gocd.Configuration{}

	qs := []*survey.Question{
		{
			Name:     "gocd_server",
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
			Name:   "ssl_check",
			Prompt: &survey.Confirm{Message: "Validate SSL Certiicate?"},
		},
	}

	a := struct {
		GoCDServer string `survey:"gocd_server"`
		Username   string
		Password   string
		SslCheck   bool `survey:"ssl_check"`
	}{}

	survey.Ask(qs, &a)

	cfg.Server = a.GoCDServer
	cfg.Username = a.Username
	cfg.Password = a.Password
	cfg.SslCheck = a.SslCheck

	s, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}

	return string(s), nil
}

func configFilePath() string {
	// @TODO Make it work for windows. Maybe...
	usr, _ := user.Current()
	return strings.Replace(ConfigDirectoryPath, "~", usr.HomeDir, 1)
}

func loadConfig() (*gocd.Configuration, error) {
	cfg := &gocd.Configuration{}

	p := configFilePath()
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		s, err := ioutil.ReadFile(configFilePath())
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(s, &cfg)
		if err != nil {
			return nil, err
		}
	}

	if server := os.Getenv(EnvVarServer); server != "" {
		cfg.Server = server
	}

	if username := os.Getenv(EnvVarUsername); username != "" {
		cfg.Username = username
	}

	if password := os.Getenv(EnvVarPassword); password != "" {
		cfg.Password = password
	}

	return cfg, nil
}

// ConfigureCommand handles the interaction between the cli flags and the action handler for configure
func ConfigureCommand() *cli.Command {
	return &cli.Command{
		Name:   ConfigureCommandName,
		Usage:  ConfigureCommandUsage,
		Action: configureAction,
	}
}
