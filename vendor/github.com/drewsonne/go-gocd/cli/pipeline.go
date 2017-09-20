package cli

import (
	"context"
	"errors"
	"github.com/urfave/cli"
)

// List of command name and descriptions
const (
	GetPipelineHistoryCommandName   = "get-pipeline-history"
	GetPipelineCommandName          = "get-pipeline"
	GetPipelineHistoryCommandUsage  = "Get Pipeline History"
	GetPipelineCommandUsage         = "Get Pipeline"
	ReleasePipelineLockCommandName  = "release-pipeline-lock"
	ReleasePipelineLockCommandUsage = "Release Pipeline Lock"
	UnpausePipelineCommandName      = "unpause-pipeline"
	UnpausePipelineCommandUsage     = "Unpause Pipeline"
	PausePipelineCommandName        = "pause-pipeline"
	PausePipelineCommandUsage       = "Pause Pipeline"
	GetPipelineStatusCommandName    = "get-pipeline-status"
	GetPipelineStatusCommandUsage   = "Get Pipeline Status"
)

// GetPipelineStatusAction handles the business logic between the command objects and the go-gocd library.
func getPipelineStatusAction(c *cli.Context) error {
	if c.String("name") == "" {
		return handleOutput(nil, nil, "GetPipelineStatus", errors.New("'--name' is missing"))
	}
	pgs, r, err := cliAgent(c).Pipelines.GetStatus(context.Background(), c.String("name"), -1)
	if err != nil {
		return handleOutput(nil, r, "GetPipelineStatus", err)
	}

	return handleOutput(pgs, r, "GetPipelineStatus", err)
}

// GetPipelineAction handles the business logic between the command objects and the go-gocd library.
func getPipelineAction(c *cli.Context) error {
	if c.String("name") == "" {
		return handleOutput(nil, nil, "GetPipeline", errors.New("'--name' is missing"))
	}
	pgs, r, err := cliAgent(c).PipelineConfigs.Get(context.Background(), c.String("name"))
	if err != nil {
		return handleOutput(nil, r, "GetPipeline", err)
	}

	return handleOutput(pgs, r, "GetPipeline", err)
}

// GetPipelineHistoryAction handles the interaction between the cli flags and the action handler for
// get-pipeline-history-action
func getPipelineHistoryAction(c *cli.Context) error {
	if c.String("name") == "" {
		return handleOutput(nil, nil, "GetPipelineHistory", errors.New("'--name' is missing"))
	}

	pgs, r, err := cliAgent(c).Pipelines.GetHistory(context.Background(), c.String("name"), -1)
	if err != nil {
		return handleOutput(nil, r, "GetPipelineHistory", err)
	}

	return handleOutput(pgs, r, "GetPipelineHistory", err)
}

// PausePipelineAction handles the business logic between the command objects and the go-gocd library.
func pausePipelineAction(c *cli.Context) error {
	if c.String("name") == "" {
		return handleOutput(nil, nil, "PausePipeline", errors.New("'--name' is missing"))
	}

	pgs, r, err := cliAgent(c).Pipelines.Pause(context.Background(), c.String("name"))
	if err != nil {
		return handleOutput(nil, r, "PausePipeline", err)
	}

	return handleOutput(pgs, r, "PausePipeline", err)
}

// UnpausePipelineAction handles the business logic between the command objects and the go-gocd library.
func unpausePipelineAction(c *cli.Context) error {
	if c.String("name") == "" {
		return handleOutput(nil, nil, "UnpausePipeline", errors.New("'--name' is missing"))
	}

	pgs, r, err := cliAgent(c).Pipelines.Unpause(context.Background(), c.String("name"))
	if err != nil {
		return handleOutput(nil, r, "UnpausePipeline", err)
	}

	return handleOutput(pgs, r, "UnpausePipeline", err)
}

// ReleasePipelineLockAction handles the business logic between the command objects and the go-gocd library.
func releasePipelineLockAction(c *cli.Context) error {
	if c.String("name") == "" {
		return handleOutput(nil, nil, "ReleasePipelinelock", errors.New("'--name' is missing"))
	}

	pgs, r, err := cliAgent(c).Pipelines.ReleaseLock(context.Background(), c.String("name"))
	if err != nil {
		return handleOutput(nil, r, "ReleasePipelinelock", err)
	}

	return handleOutput(pgs, r, "ReleasePipelinelock", err)
}

// GetPipelineStatusCommand handles the interaction between the cli flags and the action handler for
// get-pipeline
func getPipelineStatusCommand() *cli.Command {
	return &cli.Command{
		Name:     GetPipelineStatusCommandName,
		Usage:    GetPipelineStatusCommandUsage,
		Category: "Pipelines",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
		Action: getPipelineStatusAction,
	}
}

// PausePipelineCommand handles the interaction between the cli flags and the action handler for
// get-pipeline
func pausePipelineCommand() *cli.Command {
	return &cli.Command{
		Name:     PausePipelineCommandName,
		Usage:    PausePipelineCommandUsage,
		Category: "Pipelines",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
		Action: pausePipelineAction,
	}
}

// UnpausePipelineCommand handles the interaction between the cli flags and the action handler for
// get-pipeline
func unpausePipelineCommand() *cli.Command {
	return &cli.Command{
		Name:     UnpausePipelineCommandName,
		Usage:    UnpausePipelineCommandUsage,
		Category: "Pipelines",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
		Action: unpausePipelineAction,
	}
}

// ReleasePipelineLockCommand handles the interaction between the cli flags and the action handler for
// get-pipeline
func releasePipelineLockCommand() *cli.Command {
	return &cli.Command{
		Name:     ReleasePipelineLockCommandName,
		Usage:    ReleasePipelineLockCommandUsage,
		Category: "Pipelines",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
		Action: releasePipelineLockAction,
	}
}

// GetPipelineCommand handles the interaction between the cli flags and the action handler for
// get-pipeline
func getPipelineCommand() *cli.Command {
	return &cli.Command{
		Name:     GetPipelineCommandName,
		Usage:    GetPipelineCommandUsage,
		Category: "Pipelines",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
		Action: getPipelineAction,
	}
}

// GetPipelineHistoryCommand handles the interaction between the cli flags and the action handler for
// get-pipeline-history-action
func getPipelineHistoryCommand() *cli.Command {
	return &cli.Command{
		Name:     GetPipelineHistoryCommandName,
		Usage:    GetPipelineHistoryCommandUsage,
		Category: "Pipelines",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "name"},
		},
		Action: getPipelineHistoryAction,
	}
}
