package gocd

import (
	"context"
	"net/url"
	"strings"
)

// EnvironmentsService exposes calls for interacting with Environment objects in the GoCD API.
type EnvironmentsService service

// EnvironmentsResponseLinks describes the HAL _link resource for the api response object for a collection of environment
// objects
//go:generate gocd-response-links-generator -type=EnvironmentsResponseLinks,EnvironmentLinks
type EnvironmentsResponseLinks struct {
	Self *url.URL `json:"self"`
	Doc  *url.URL `json:"doc"`
}

// EnvironmentLinks describes the HAL _link resource for the api response object for a collection of environment objects.
type EnvironmentLinks struct {
	Self *url.URL `json:"self"`
	Doc  *url.URL `json:"doc"`
	Find *url.URL `json:"find"`
}

// EnvironmentsResponse describes the response obejct for a plugin API call.
type EnvironmentsResponse struct {
	Links    *EnvironmentsResponseLinks `json:"_links"`
	Embedded struct {
		Environments []*Environment `json:"environments"`
	} `json:"_embedded"`
}

// Environment describes a group of pipelines and agents
type Environment struct {
	Links                *EnvironmentLinks      `json:"_links,omitempty"`
	Name                 string                 `json:"name"`
	Pipelines            []*Pipeline            `json:"pipelines,omitempty"`
	Agents               []*Agent               `json:"agents,omitempty"`
	EnvironmentVariables []*EnvironmentVariable `json:"environment_variables,omitempty"`
	Version              string                 `json:"version"`
}

// EnvironmentPatchRequest describes the actions to perform on an environment
type EnvironmentPatchRequest struct {
	Pipelines            *PatchStringAction          `json:"pipelines"`
	Agents               *PatchStringAction          `json:"agents"`
	EnvironmentVariables *EnvironmentVariablesAction `json:"environment_variables"`
}

// EnvironmentVariablesAction describes a collection of Environment Variables to add or remove.
type EnvironmentVariablesAction struct {
	Add    []*EnvironmentVariable `json:"add"`
	Remove []*EnvironmentVariable `json:"remove"`
}

// PatchStringAction describes a collection of resources to add or remove.
type PatchStringAction struct {
	Add    []string `json:"add"`
	Remove []string `json:"remove"`
}

// List all environments
func (es *EnvironmentsService) List(ctx context.Context) (*EnvironmentsResponse, *APIResponse, error) {
	e := EnvironmentsResponse{}
	_, resp, err := es.client.getAction(ctx, &APIClientRequest{
		Path:         "admin/environments",
		ResponseBody: &e,
		APIVersion:   apiV2,
	})

	return &e, resp, err
}

// Get a single environment by name
func (es *EnvironmentsService) Get(ctx context.Context, name string) (*Environment, *APIResponse, error) {
	e := Environment{}
	_, resp, err := es.client.getAction(ctx, &APIClientRequest{
		Path:         "admin/environments/" + name,
		ResponseBody: &e,
		APIVersion:   apiV2,
	})
	if err == nil {
		e.Version = strings.Replace(resp.HTTP.Header.Get("Etag"), "\"", "", -1)
	}

	return &e, resp, err
}

// Patch an environments configuration by adding or removing pipelines, agents, environment variables
func (es *EnvironmentsService) Patch(ctx context.Context, name string, patch *EnvironmentPatchRequest) (*Environment, *APIResponse, error) {
	env := Environment{}
	_, resp, err := es.client.patchAction(ctx, &APIClientRequest{
		Path:         "admin/environments" + name,
		RequestBody:  patch,
		ResponseBody: &env,
		APIVersion:   apiV2,
	})
	if err == nil {
		env.Version = strings.Replace(resp.HTTP.Header.Get("Etag"), "\"", "", -1)
	}

	return &env, resp, err
}
