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
	pType, pipeline, name, err := parseGoCDPipelineStageId(d)
	if err != nil {
		return nil, err
	}

	d.Set("name", name)

	if pType == STAGE_TYPE_PIPELINE {
		d.Set("pipeline", pipeline)
	} else if pType == STAGE_TYPE_PIPELINE_TEMPLATE {
		d.Set("pipeline_template", pipeline)
	} else {
		return nil, fmt.Errorf("Unexpected pipeline type `%s`", pType)
	}

	if err := resourcePipelineStageRead(d, meta); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func resourcePipelineStageExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	pType, pipeline, name, err := parseGoCDPipelineStageId(d)

	stage, err := retrieveStage(pType, name, pipeline, d, meta)
	if err != nil {
		return false, err
	}

	if stage == nil && err == nil {
		return false, nil
	}
	return true, nil
}

func resourcePipelineStageCreate(d *schema.ResourceData, meta interface{}) error {
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
	ctx := context.Background()
	client.Lock()
	defer client.Unlock()

	if pipelineTemplateI, hasPipelineTemplate := d.GetOk("pipeline_template"); hasPipelineTemplate {
		pipelineTemplate := pipelineTemplateI.(string)
		existingPt, _, err := client.PipelineTemplates.Get(ctx, pipelineTemplate)
		if err != nil {
			return err
		}

		existingPt.Stages = cleanPlaceHolderStage(existingPt.Stages)
		existingPt.Stages = append(existingPt.Stages, &doc)

		pt, _, err := client.PipelineTemplates.Update(ctx, pipelineTemplate, existingPt)
		if err != nil {
			return err
		}

		d.SetId(fmt.Sprintf("template/%s/%s", pt.Name, doc.Name))
	} else if pipelineI, hasPipeline := d.GetOk("pipeline"); hasPipeline {
		pipeline := pipelineI.(string)
		existingP, _, err := client.PipelineConfigs.Get(ctx, pipeline)
		if err != nil {
			return err
		}

		existingP.Stages = cleanPlaceHolderStage(existingP.Stages)
		existingP.Stages = append(existingP.Stages, &doc)

		p, _, err := client.PipelineConfigs.Update(ctx, existingP.Name, existingP)

		d.SetId(fmt.Sprintf("pipeline/%s/%s", p.Name, doc.Name))
	}

	return nil
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

	if pType == STAGE_TYPE_PIPELINE {
		d.Set("pipeline", pipeline)
	} else if pType == STAGE_TYPE_PIPELINE_TEMPLATE {
		d.Set("pipeline_template", pipeline)
	} else {
		return fmt.Errorf("Unexpected pipeline type `%s`", pType)
	}

	return nil
}

func resourcePipelineStageUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePipelineStageDelete(d *schema.ResourceData, meta interface{}) error {
	var stageName string
	if name, hasName := d.GetOk("name"); hasName {
		stageName = name.(string)
	} else {
		return errors.New("Missing `name`")
	}

	name, pType, err := pipelineNameType(d)
	if err != nil {
		return err
	}

	client := meta.(*gocd.Client)
	ctx := context.Background()
	var existing gocd.StageContainer

	// Make delete operation atomic
	client.Lock()
	defer client.Unlock()

	if pType == STAGE_TYPE_PIPELINE_TEMPLATE {
		existing, _, err = client.PipelineTemplates.Get(ctx, name)
	} else if pType == STAGE_TYPE_PIPELINE {
		existing, _, err = client.PipelineConfigs.Get(ctx, name)
	}
	if err != nil {
		return err
	}

	stages := cleanPlaceHolderStage(existing.GetStages())

	cleanedStages := []*gocd.Stage{}
	for _, stage := range stages {
		if stage.Name != stageName {
			cleanedStages = append(cleanedStages, stage)
		}
	}

	if len(cleanedStages) == 0 {
		cleanedStages = append(cleanedStages, stagePlaceHolder())
	}

	existing.SetStages(cleanedStages)

	var updated gocd.StageContainer
	if pType == STAGE_TYPE_PIPELINE_TEMPLATE {
		updated, _, err = client.PipelineTemplates.Update(ctx, existing.GetName(), existing.(*gocd.PipelineTemplate))
	} else if pType == STAGE_TYPE_PIPELINE {
		updated, _, err = client.PipelineConfigs.Update(ctx, existing.GetName(), existing.(*gocd.Pipeline))
	}

	if err != nil {
		return err
	}

	stageExists := false
	for _, stage := range updated.GetStages() {
		if stage.Name == stageName {
			stageExists = true
		}
	}

	if stageExists {
		return fmt.Errorf("Could not delete stage `%s`", stageName)
	}

	return nil
}

func retrieveStage(pType string, stageName string, pipeline string, d *schema.ResourceData, meta interface{}) (*gocd.Stage, error) {

	var stages []*gocd.Stage
	var existing gocd.StageContainer
	var err error

	client := meta.(*gocd.Client)
	ctx := context.Background()

	if pType == STAGE_TYPE_PIPELINE_TEMPLATE {
		existing, _, err = client.PipelineTemplates.Get(ctx, pipeline)
	} else if pType == STAGE_TYPE_PIPELINE {
		existing, _, err = client.PipelineConfigs.Get(ctx, pipeline)
	}
	if err != nil {
		return nil, err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", pType, existing.GetName(), stageName))
	stages = cleanPlaceHolderStage(existing.GetStages())

	for _, stage := range stages {
		if stage.Name == stageName {
			return stage, nil
		}
	}

	return nil, nil
}

func pipelineNameType(d *schema.ResourceData) (name string, pType string, err error) {
	if pipelineTemplateI, hasPipelineTemplate := d.GetOk("pipeline_template"); hasPipelineTemplate {
		return pipelineTemplateI.(string), STAGE_TYPE_PIPELINE_TEMPLATE, nil
	} else if pipelineI, hasPipeline := d.GetOk("pipeline"); hasPipeline {
		return pipelineI.(string), STAGE_TYPE_PIPELINE, nil
	}
	return "", "", errors.New("Could not find `pipeline` nor `pipeline_template`")
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
	r, err := regexp.Compile(`^(template|pipeline)/([^/]+)/([^/]+)$`)
	if err != nil {
		return "", "", "", err
	}

	id := d.Id()
	matches := r.FindAllStringSubmatch(id, -1)

	if len(matches) != 1 {
		return "", "", "", fmt.Errorf("Could not parse the provided id `%s`", id)
	}

	return matches[0][1], matches[0][2], matches[0][3], nil
}
