package cli

import (
	"context"
	"errors"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
	"strings"
)

// List of command name and descriptions
const (
	ListEnvironmentsCommandName                = "list-environments"
	ListEnvironmentsCommandUsage               = "List all environments"
	GetEnvironmentCommandName                  = "get-environment"
	GetEnvironmentCommandUsage                 = "Get an environment by name"
	AddPipelinesToEnvironmentCommandName       = "add-pipelines-to-environment"
	AddPipelinesToEnvironmentCommandUsage      = "Add one or more pipelines to an environment"
	RemovePipelinesFromEnvironmentCommandName  = "remove-pipelines-from-environment"
	RemovePipelinesFromEnvironmentCommandUsage = "Remove one or more pipelines from an environment"
)

// ListEnvironmentsAction handles the listing of environments
func listEnvironmentsAction(c *cli.Context) error {
	es, r, err := cliAgent(c).Environments.List(context.Background())
	if err != nil {
		return handleErrOutput("ListEnvironments", err)
	}

	es.RemoveLinks()

	return handleOutput(es, r, "ListEnvironments", err)
}

// GetEnvironmentAction handles the retrieval of environments
func getEnvironmentAction(c *cli.Context) error {
	if c.String("name") == "" {
		return handleOutput(nil, nil, "GetEnvironment", errors.New("'--name' is missing"))
	}
	e, r, err := cliAgent(c).Environments.Get(context.Background(), c.String("name"))
	if err != nil {
		return handleErrOutput("GetEnvironment", err)
	}
	e.RemoveLinks()
	return handleOutput(e, r, "GetEnvironment", err)
}

// AddPipelinesToEnvironmentAction handles the adding of a pipeline to an environment
func addPipelinesToEnvironmentAction(c *cli.Context) error {
	if c.String("environment-name") == "" {
		return handleOutput(nil, nil, "AddPipelinesToEnvironment", errors.New("'--environment-name' is missing"))
	}
	if c.String("pipeline-names") == "" {
		return handleOutput(nil, nil, "AddPipelinesToEnvironment", errors.New("'--pipeline-names' is missing"))
	}

	e, r, err := cliAgent(c).Environments.Patch(context.Background(), c.String("environment-name"), &gocd.EnvironmentPatchRequest{
		Pipelines: &gocd.PatchStringAction{
			Add: strings.Split(c.String("pipeline-names"), ","),
		},
	})
	if err != nil {
		return handleErrOutput("AddPipelinesToEnvironment", err)
	}
	e.RemoveLinks()
	return handleOutput(e, r, "AddPipelinesToEnvironment", err)
}

// RemovePipelinesFromEnvironmentAction handles the removing of a pipeline from an environment
func removePipelinesFromEnvironmentAction(c *cli.Context) error {
	if c.String("environment-name") == "" {
		return handleOutput(nil, nil, "RemovePipelinesFromEnvironment", errors.New("'--environment-name' is missing"))
	}
	if c.String("pipeline-names") == "" {
		return handleOutput(nil, nil, "RemovePipelinesFromEnvironment", errors.New("'--pipeline-names' is missing"))
	}

	e, r, err := cliAgent(c).Environments.Patch(context.Background(), c.String("environment-name"), &gocd.EnvironmentPatchRequest{
		Pipelines: &gocd.PatchStringAction{
			Remove: strings.Split(c.String("pipeline-names"), ","),
		},
	})
	if err != nil {
		return handleErrOutput("RemovePipelinesFromEnvironment", err)
	}
	e.RemoveLinks()
	return handleOutput(e, r, "RemovePipelinesFromEnvironment", err)
}

// ListEnvironmentsCommand handles definition of cli command
func listEnvironmentsCommand() *cli.Command {
	return &cli.Command{
		Name:     ListEnvironmentsCommandName,
		Usage:    ListEnvironmentsCommandUsage,
		Action:   listEnvironmentsAction,
		Category: "Environments",
	}
}

// GetEnvironmentCommand handles definition of cli command
func getEnvironmentCommand() *cli.Command {
	return &cli.Command{
		Name:     GetEnvironmentCommandName,
		Usage:    GetEnvironmentCommandUsage,
		Action:   getEnvironmentAction,
		Category: "Environments",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
	}
}

// AddPipelinesToEnvironmentCommand handles definition of cli command
func addPipelinesToEnvironmentCommand() *cli.Command {
	return &cli.Command{
		Name:     AddPipelinesToEnvironmentCommandName,
		Usage:    AddPipelinesToEnvironmentCommandUsage,
		Action:   addPipelinesToEnvironmentAction,
		Category: "Environments",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "environment-name"},
			cli.StringFlag{
				Name:  "pipeline-names",
				Usage: "Comma seperated list of pipeline names to add.",
			},
		},
	}
}

// RemovePipelinesFromEnvironmentCommand handles definition of cli command
func removePipelinesFromEnvironmentCommand() *cli.Command {
	return &cli.Command{
		Name:     RemovePipelinesFromEnvironmentCommandName,
		Usage:    RemovePipelinesFromEnvironmentCommandUsage,
		Action:   removePipelinesFromEnvironmentAction,
		Category: "Environments",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "environment-name"},
			cli.StringFlag{
				Name:  "pipeline-names",
				Usage: "Comma seperated list of pipeline names to remove.",
			},
		},
	}
}
