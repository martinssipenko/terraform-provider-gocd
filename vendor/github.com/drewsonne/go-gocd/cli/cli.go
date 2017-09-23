package cli

import (
	"encoding/json"
	"fmt"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/urfave/cli"
)

// GetCliCommands returns a list of all CLI Command structs
func GetCliCommands() []cli.Command {
	return []cli.Command{
		*configureCommand(),
		*listAgentsCommand(),
		*listPipelineTemplatesCommand(),
		*getAgentCommand(),
		*getPipelineTemplateCommand(),
		*createPipelineTemplateCommand(),
		*updateAgentCommand(),
		*updateAgentsCommand(),
		*updatePipelineConfigCommand(),
		*updatePipelineTemplateCommand(),
		*deleteAgentCommand(),
		*deleteAgentsCommand(),
		*deletePipelineTemplateCommand(),
		*deletePipelineConfigCommand(),
		*listPipelineGroupsCommand(),
		*getPipelineHistoryCommand(),
		*getPipelineCommand(),
		*createPipelineConfigCommand(),
		*generateJSONSchemaCommand(),
		*getPipelineStatusCommand(),
		*pausePipelineCommand(),
		*unpausePipelineCommand(),
		*releasePipelineLockCommand(),
		*getConfigurationCommand(),
		*encryptCommand(),
		*getVersionCommand(),
		*listPluginsCommand(),
		*getPluginCommand(),
		*listScheduledJobsCommand(),
		*getPipelineConfigCommand(),
		*listEnvironmentsCommand(),
		*getEnvironmentCommand(),
		*addPipelinesToEnvironmentCommand(),
		*removePipelinesFromEnvironmentCommand(),
	}
}

// NewCliClient
func cliAgent(c *cli.Context) *gocd.Client {
	var profile string
	if profile = c.String("profile"); profile == "" {
		profile = "default"
	}

	cfg, err := gocd.LoadConfigByName(profile)
	if err != nil {
		panic(err)
	}

	if server := c.String("server"); server != "" {
		cfg.Server = server
	}

	if username := c.String("username"); username != "" {
		cfg.Username = username
	}

	if password := c.String("password"); password != "" {
		cfg.Password = password
	}

	cfg.SkipSslCheck = cfg.SkipSslCheck || c.Bool("skip_ssl_check")

	return cfg.Client()
}

func handleErrOutput(reqType string, err error) error {
	return handleOutput(nil, nil, reqType, err)
}

func handleOutput(r interface{}, hr *gocd.APIResponse, reqType string, err error) error {
	var b []byte
	var o map[string]interface{}
	if err != nil {
		o = map[string]interface{}{
			"Request": reqType,
			"Error":   err.Error(),
		}
	} else if hr.HTTP.StatusCode >= 200 && hr.HTTP.StatusCode < 300 {
		o = map[string]interface{}{
			fmt.Sprintf("%sResponse", reqType): r,
		}
		//} else if hr.HTTP.StatusCode == 404 {
		//	o = map[string]interface{}{
		//		"Error": fmt.Sprintf("Could not find resource for '%s' action.", reqType),
		//	}
	} else {

		b1, _ := json.Marshal(hr.HTTP.Header)
		b2, _ := json.Marshal(hr.Request.HTTP.Header)
		o = map[string]interface{}{
			"Error":           "An error occurred while retrieving the resource.",
			"Status":          hr.HTTP.StatusCode,
			"ResponseHeader":  string(b1),
			"ResponseBody":    hr.Body,
			"RequestBody":     hr.Request.Body,
			"RequestEndpoint": hr.Request.HTTP.URL.String(),
			"RequestHeader":   string(b2),
		}
	}
	b, err = json.MarshalIndent(o, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

	return nil
}
