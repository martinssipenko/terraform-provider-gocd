package gocd

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// PipelineTemplatesService describes the HAL _link resource for the api response object for a pipeline configuration objects.
type PipelineTemplatesService service

// PipelineTemplatesLinks describes a single pipeline template config object HAL links
//go:generate gocd-response-links-generator -type=PipelineTemplatesLinks,PipelineTemplateLinks
type PipelineTemplatesLinks struct {
	Self *url.URL `json:"self"`
	Doc  *url.URL `json:"doc"`
	Find *url.URL `json:"find"`
}

// PipelineTemplateLinks describes multiple pipeline template config object HAL links
type PipelineTemplateLinks struct {
	Self *url.URL `json:"self"`
	Doc  *url.URL `json:"doc"`
	Find *url.URL `json:"find"`
}

// PipelineTemplateRequest describes a PipelineTemplate
type PipelineTemplateRequest struct {
	Name   string   `json:"name"`
	Stages []*Stage `json:"stages"`
}

// PipelineTemplateResponse describes an api response for a single pipeline templates
type PipelineTemplateResponse struct {
	Name     string `json:"name"`
	Embedded *struct {
		Pipelines []*struct {
			Name string `json:"name"`
		}
	} `json:"_embedded,omitempty"`
}

// PipelineTemplatesResponse describes an api response for multiple pipeline templates
type PipelineTemplatesResponse struct {
	Links    PipelineTemplatesLinks `json:"_links,omitempty"`
	Embedded *struct {
		Templates []*PipelineTemplate `json:"templates"`
	} `json:"_embedded,omitempty"`
}

// PipelineTemplate describes a response from the API for a pipeline template object.
type PipelineTemplate struct {
	Links    *PipelineTemplateLinks `json:"_links,omitempty"`
	Name     string                 `json:"name"`
	Embedded *struct {
		Pipelines []*Pipeline `json:"pipelines,omitempty"`
	} `json:"_embedded,omitempty"`
	Version string   `json:"template_version"`
	Stages  []*Stage `json:"stages,omitempty"`
}

// RemoveLinks gets the PipelineTemplate ready to be submitted to the GoCD API.
func (pt *PipelineTemplate) RemoveLinks() {
	pt.Links = nil
}

// Pipelines returns a list of Pipelines attached to this PipelineTemplate object.
func (pt *PipelineTemplate) Pipelines() []*Pipeline {
	return pt.Embedded.Pipelines
}

// Exists ensures whether or not a PipelineTemplate is present in the API.
func (pts *PipelineTemplatesService) Exists(ctx context.Context, name string) (bool, *APIResponse, error) {
	return pts.client.genericHeadAction(ctx, fmt.Sprintf("admin/templates/%s", name), apiV3)
}

// Get a single PipelineTemplate object in the GoCD API.
func (pts *PipelineTemplatesService) Get(ctx context.Context, name string) (*PipelineTemplate, *APIResponse, error) {
	pt := PipelineTemplate{}
	_, resp, err := pts.client.getAction(ctx, &APIClientRequest{
		Path:         "admin/templates/" + name,
		APIVersion:   apiV3,
		ResponseBody: &pt,
	})
	return &pt, resp, err
}

// List all PipelineTemplate objects in the GoCD API.
func (pts *PipelineTemplatesService) List(ctx context.Context) ([]*PipelineTemplate, *APIResponse, error) {
	ptr := PipelineTemplatesResponse{}

	_, resp, err := pts.client.getAction(ctx, &APIClientRequest{
		Path:         "admin/templates",
		APIVersion:   apiV3,
		ResponseBody: &ptr,
	})

	return ptr.Embedded.Templates, resp, err
}

// Create a new PipelineTemplate object in the GoCD API.
func (pts *PipelineTemplatesService) Create(ctx context.Context, name string, st []*Stage) (*PipelineTemplate, *APIResponse, error) {

	pt := PipelineTemplateRequest{
		Name:   name,
		Stages: st,
	}
	ptr := PipelineTemplate{}

	_, resp, err := pts.client.putAction(ctx, &APIClientRequest{
		Path:         "admin/templates",
		APIVersion:   apiV3,
		RequestBody:  pt,
		ResponseBody: &ptr,
	})

	ptr.Version = strings.Replace(resp.HTTP.Header.Get("Etag"), "\"", "", -1)

	return &ptr, resp, err

}

// Update an PipelineTemplate object in the GoCD API.
func (pts *PipelineTemplatesService) Update(ctx context.Context, name string, version string, st []*Stage) (*PipelineTemplate, *APIResponse, error) {
	pt := &PipelineTemplateRequest{
		Name:   name,
		Stages: st,
	}
	ptr := &PipelineTemplate{}

	_, resp, err := pts.client.putAction(ctx, &APIClientRequest{
		Path:         "admin/templates/" + name,
		APIVersion:   apiV3,
		RequestBody:  pt,
		ResponseBody: &ptr,
	})

	return ptr, resp, err

}

// Delete a PipelineTemplate from the GoCD API.
func (pts *PipelineTemplatesService) Delete(ctx context.Context, name string) (string, *APIResponse, error) {
	return pts.client.deleteAction(ctx, fmt.Sprintf("admin/templates/%s", name), apiV3)
}
