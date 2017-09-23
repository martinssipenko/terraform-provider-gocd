package cli

import (
	"context"
	"github.com/urfave/cli"
)

// List of command name and descriptions
const (
	GetConfigurationCommandName  = "get-configuration"
	GetConfigurationCommandUsage = "Get GoCD server configuration. This is the cruise-config.xml file. It is exposed in a json format to enable a consistent format. This API is for read-only purposes and not intended as an interface to modify the config."
	GetVersionCommandName        = "get-version"
	GetVersionCommandUsage       = "Return the Version for the GoCD Server"
)

// GetConfigurationAction gets a list of agents and return them.
func getConfigurationAction(c *cli.Context) (err error) {
	if pgs, r, err := cliAgent(c).Configuration.Get(context.Background()); err == nil {
		return handleOutput(pgs, r, "GetConfiguration", err)
	}
	return handleErrOutput("GetConfiguration", err)
}

// GetVersionAction returns version information about GoCD
func getVersionAction(c *cli.Context) (err error) {
	if v, r, err := cliAgent(c).Configuration.GetVersion(context.Background()); err == nil {
		return handleOutput(v, r, "GetVersion", err)
	}
	return handleErrOutput("GetVersion", err)
}

// GetConfigurationCommand handles the interaction between the cli flags and the action handler for delete-agents
func getConfigurationCommand() *cli.Command {
	return &cli.Command{
		Name:     GetConfigurationCommandName,
		Usage:    GetConfigurationCommandUsage,
		Action:   getConfigurationAction,
		Category: "Configuration",
	}
}

// GetVersionCommand handles the interaction between the cli flags and the action handler for delete-agents
func getVersionCommand() *cli.Command {
	return &cli.Command{
		Name:     GetVersionCommandName,
		Usage:    GetVersionCommandUsage,
		Action:   getVersionAction,
		Category: "Configuration",
	}
}
