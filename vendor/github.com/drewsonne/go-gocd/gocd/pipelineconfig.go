package gocd

import (
	"context"
	"fmt"
	"strings"
)

// PipelineConfigsService describes the HAL _link resource for the api response object for a pipelineconfig
type PipelineConfigsService service

// PipelineConfig describes the configuration for a pipeline
type PipelineConfig struct{}

// PipelineConfigRequest describes a request object for creating or updating pipelines
type PipelineConfigRequest struct {
	Group    string    `json:"group"`
	Pipeline *Pipeline `json:"pipeline"`
}

// Update a pipeline configuration
func (pcs *PipelineConfigsService) Update(ctx context.Context, group string, name string, version string, p *Pipeline) (*Pipeline, *APIResponse, error) {

	pt := &PipelineConfigRequest{
		Group:    group,
		Pipeline: p,
	}

	req, err := pcs.client.NewRequest("PUT", "admin/pipelines/"+name, pt, apiV4)
	if err != nil {
		return nil, nil, err
	}

	req.HTTP.Header.Set("If-Match", fmt.Sprintf("\"%s\"", version))

	pc := Pipeline{}
	resp, err := pcs.client.Do(ctx, req, &pc, responseTypeJSON)
	if err != nil {
		return nil, resp, err
	}

	return &pc, resp, nil

}

// Create a pipeline configuration
func (pcs *PipelineConfigsService) Create(ctx context.Context, group string, p *Pipeline) (*Pipeline, *APIResponse, error) {
	pt := &PipelineConfigRequest{
		Group:    group,
		Pipeline: p,
	}

	req, err := pcs.client.NewRequest("POST", "admin/pipelines", pt, apiV4)
	if err != nil {
		return nil, nil, err
	}

	pc := Pipeline{}
	resp, err := pcs.client.Do(ctx, req, &pc, responseTypeJSON)
	if err != nil {
		return nil, resp, err
	}
	pc.Version = strings.Replace(resp.HTTP.Header.Get("Etag"), "\"", "", -1)

	return &pc, resp, nil
}

// Delete a pipeline configuration
func (pcs *PipelineConfigsService) Delete(ctx context.Context, name string) (string, *APIResponse, error) {
	return pcs.client.deleteAction(ctx, "admin/pipelines/"+name, apiV3)
}
