package gocd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
	"regexp"
)

const STAGE_TYPE_PIPELINE = "pipeline"
const STAGE_TYPE_PIPELINE_TEMPLATE = "template"

func resourcePipelineStage() *schema.Resource {
	stringArg := &schema.Schema{Type: schema.TypeString}
	optionalBoolArg := &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}

	return &schema.Resource{
		Create: resourcePipelineStageCreate,
		Read:   resourcePipelineStageRead,
		Update: resourcePipelineStageUpdate,
		Delete: resourcePipelineStageDelete,
		Exists: resourcePipelineStageExists,
		Importer: &schema.ResourceImporter{
			State: resourcePipelineStageImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fetch_materials":         optionalBoolArg,
			"clean_working_directory": optionalBoolArg,
			"never_cleanup_artifacts": optionalBoolArg,
			"jobs": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     stringArg,
			},
			"manual_approval": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"success_approval"},
			},
			"success_approval": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"manual_approval"},
			},
			"authorization_users": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"success_approval", "authorization_roles"},
				Elem:          stringArg,
			},
			"authorization_roles": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"success_approval", "authorization_users"},
				Elem:          stringArg,
			},
			"environment_variables": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     stringArg,
			},
			"pipeline": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"pipeline_template"},
				Optional:      true,
			},
			"pipeline_template": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"pipeline"},
				Optional:      true,
			},
		},
	}
}

func resourcePipelineStageImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var pType, pipeline, name string
	var err error

	if pType, pipeline, name, err = parseGoCDPipelineStageId(d); err != nil {
		return nil, err
	}

	d.Set("name", name)

	if err := resourcePipelineStageSetPTypeName(d, pType, pipeline); err != nil {
		return nil, err
	}

	if err := resourcePipelineStageRead(d, meta); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func resourcePipelineStageExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	var pType, pipeline, name string
	var err error

	if pType, pipeline, name, err = parseGoCDPipelineStageId(d); err != nil {
		return false, err
	}

	if stage, err := retrieveStage(pType, name, pipeline, d, meta); err == nil {
		return stage != nil, nil
	} else {
		return false, err
	}

}

func resourcePipelineStageCreate(d *schema.ResourceData, meta interface{}) error {
	var existing *gocd.StageContainer
	var updated *gocd.StageContainer

	var pipelineName, pType string

	var err error

	doc := gocd.Stage{Approval: &gocd.Approval{}}

	if name, hasName := d.GetOk("name"); hasName {
		doc.Name = name.(string)
	} else {
		return errors.New("Missing `name`")
	}

	if manualApproval, ok := d.GetOk("manual_approval"); ok && manualApproval.(bool) {
		dataSourceStageParseManuallApproval(d, &doc)
	} else if d.Get("success_approval").(bool) {
		doc.Approval.Type = "success"
		doc.Approval.Authorization = nil
	}

	if rJobs, hasJobs := d.GetOk("jobs"); hasJobs {
		if jobs := decodeConfigStringList(rJobs.([]interface{})); len(jobs) > 0 {
			dataSourceStageParseJobs(jobs, &doc)
		}
	}

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	if pipelineName, pType, err = pipelineNameType(d); err != nil {
		return err
	}

	existing, err = getStageContainer(pType, pipelineName, client)
	if err != nil {
		return err
	}

	(*existing).SetStages(cleanPlaceHolderStage((*existing).GetStages()))
	(*existing).AddStage(&doc)

	if updated, err = updateStageContainer(pType, existing, client); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", pType, (*updated).GetName(), doc.Name))

	return err
}

func resourcePipelineStageRead(d *schema.ResourceData, meta interface{}) error {
	pType, pipeline, name, err := parseGoCDPipelineStageId(d)
	stage, err := retrieveStage(pType, name, pipeline, d, meta)
	if err != nil {
		return err
	}

	d.Set("name", stage.Name)
	d.Set("fetch_materials", stage.FetchMaterials)
	d.Set("clean_working_directory", stage.CleanWorkingDirectory)
	d.Set("never_cleanup_artifacts", stage.NeverCleanupArtifacts)
	d.Set("environment_variables", stage.EnvironmentVariables)
	d.Set("resources", stage.Resources)

	if appr := stage.Approval; appr != nil {
		if appr.Type == "manual" {
			d.Set("manual_approval", true)
		} else if appr.Type == "success_approval" {
			d.Set("success_approval", true)
		}
		if auth := stage.Approval.Authorization; auth != nil {
			if users := auth.Users; users != nil {
				d.Set("authorization_users", users)
			}
			if roles := auth.Roles; roles != nil {
				d.Set("authorization_roles", roles)
			}
		}
	}

	if jobs := stage.Jobs; len(jobs) > 0 {
		stringJobs := []string{}
		for _, job := range jobs {
			s, err := job.JSONString()
			if err != nil {
				return err
			}
			stringJobs = append(stringJobs, s)
		}

		d.Set("jobs", stringJobs)
	}

	if err := resourcePipelineStageSetPTypeName(d, pType, pipeline); err != nil {
		return err
	}

	return nil
}

func resourcePipelineStageUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePipelineStageDelete(d *schema.ResourceData, meta interface{}) error {
	var updated *gocd.StageContainer
	var existing *gocd.StageContainer
	var pipeline, pType string
	var err error

	if pipeline, pType, err = pipelineNameType(d); err != nil {
		return err
	}

	client := meta.(*gocd.Client)

	// Make delete operation atomic
	client.Lock()
	defer client.Unlock()

	// Retrieve Pipeline object so we have the latest version
	if existing, err = getStageContainer(pType, pipeline, client); err != nil {
		return err
	}

	stages := cleanPlaceHolderStage((*existing).GetStages())

	cleanedStages := []*gocd.Stage{}
	stageName := d.Get("name").(string)
	for _, stage := range stages {
		if stage.Name != stageName {
			cleanedStages = append(cleanedStages, stage)
		}
	}

	if len(cleanedStages) == 0 {
		cleanedStages = append(cleanedStages, stagePlaceHolder())
	}

	(*existing).SetStages(cleanedStages)

	// Perform stage update
	if updated, err = updateStageContainer(pType, existing, meta); err != nil {
		return err
	}

	if stage := (*updated).GetStage(stageName); stage != nil {
		return fmt.Errorf("Could not delete stage `%s`. Does not exist.", stageName)
	}

	return nil
}

func retrieveStage(pType string, stageName string, pipeline string, d *schema.ResourceData, meta interface{}) (*gocd.Stage, error) {
	var existing *gocd.StageContainer
	var err error

	client := meta.(*gocd.Client)

	if existing, err = getStageContainer(pType, pipeline, client); err != nil {
		return nil, err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", pType, (*existing).GetName(), stageName))

	if stage := (*existing).GetStage(stageName); stage != nil {
		return stage, nil
	}

	return nil, nil
}

func pipelineNameType(d *schema.ResourceData) (pipelineName string, pType string, err error) {
	if pipelineTemplateI, hasPipelineTemplate := d.GetOk("pipeline_template"); hasPipelineTemplate {
		return pipelineTemplateI.(string), STAGE_TYPE_PIPELINE_TEMPLATE, nil
	} else if pipelineI, hasPipeline := d.GetOk("pipeline"); hasPipeline {
		return pipelineI.(string), STAGE_TYPE_PIPELINE, nil
	}
	return "", "", errors.New("Could not find `pipeline` nor `pipeline_template`")
}

func updateStageContainer(pType string, existing *gocd.StageContainer, meta interface{}) (*gocd.StageContainer, error) {
	var updated gocd.StageContainer
	var err error
	client := meta.(*gocd.Client)
	ctx := context.Background()
	if pType == STAGE_TYPE_PIPELINE_TEMPLATE {
		updated, _, err = client.PipelineTemplates.Update(ctx, (*existing).GetName(), (*existing).(*gocd.PipelineTemplate))
	} else if pType == STAGE_TYPE_PIPELINE {
		updated, _, err = client.PipelineConfigs.Update(ctx, (*existing).GetName(), (*existing).(*gocd.Pipeline))
	}
	return &updated, err
}

func getStageContainer(pType string, pipelineName string, client *gocd.Client) (*gocd.StageContainer, error) {
	var existing gocd.StageContainer
	var err error
	ctx := context.Background()
	if pType == STAGE_TYPE_PIPELINE_TEMPLATE {
		existing, _, err = client.PipelineTemplates.Get(ctx, pipelineName)
	} else if pType == STAGE_TYPE_PIPELINE {
		existing, _, err = client.PipelineConfigs.Get(ctx, pipelineName)
	}

	return &existing, err

}

func cleanPlaceHolderStage(stages []*gocd.Stage) []*gocd.Stage {
	cleanStages := []*gocd.Stage{}
	for _, stage := range stages {
		if stage.Name != PLACEHOLDER_NAME {
			cleanStages = append(cleanStages, stage)
		}
	}
	return cleanStages
}

func stagePlaceHolder() *gocd.Stage {
	return &gocd.Stage{
		Name: PLACEHOLDER_NAME,
		Jobs: []*gocd.Job{
			{Name: PLACEHOLDER_NAME},
		},
	}
}

func dataSourceStageParseManuallApproval(data *schema.ResourceData, doc *gocd.Stage) error {
	doc.Approval.Type = "manual"
	doc.Approval.Authorization = &gocd.Authorization{}
	if users := data.Get("authorization_users").(*schema.Set).List(); len(users) > 0 {
		doc.Approval.Authorization.Users = decodeConfigStringList(users)
	} else if roles := data.Get("authorization_roles").(*schema.Set).List(); len(roles) > 0 {
		doc.Approval.Authorization.Roles = decodeConfigStringList(roles)
	}
	return nil
}

func dataSourceStageParseJobs(jobs []string, doc *gocd.Stage) error {
	for _, rawjob := range jobs {
		job := gocd.Job{}
		if err := json.Unmarshal([]byte(rawjob), &job); err != nil {
			return err
		}
		doc.Jobs = append(doc.Jobs, &job)
	}
	return nil
}

func parseGoCDPipelineStageId(d *schema.ResourceData) (pType string, pipeline string, stage string, err error) {
	var r *regexp.Regexp
	var matches [][]string

	if r, err = regexp.Compile(`^(template|pipeline)/([^/]+)/([^/]+)$`); err != nil {
		return "", "", "", err
	}

	id := d.Id()
	if matches = r.FindAllStringSubmatch(id, -1); len(matches) != 1 {
		return "", "", "", fmt.Errorf("Could not parse the provided id `%s`", id)
	}

	return matches[0][1], matches[0][2], matches[0][3], nil
}
