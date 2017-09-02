package gocd

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/drewsonne/go-gocd/gocd"
	"errors"
	"fmt"
	"context"
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
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("name", d.Id())
				return []*schema.ResourceData{d}, nil
			},
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

func resourcePipelineStageExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	stage, err := retrieveStage(d, meta)
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

		pt, _, err := client.PipelineTemplates.Update(ctx, pipelineTemplate, existingPt.Version, existingPt.Stages)
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
	stage, err := retrieveStage(d, meta)
	if err != nil {
		return err
	}

	d.Set("name", stage.Name)
	d.Set("fetch_materials", stage.FetchMaterials)
	d.Set("clean_working_directory", stage.CleanWorkingDirectory)
	d.Set("never_cleanup_artifacts", stage.NeverCleanupArtifacts)
	d.Set("environment_variables", stage.EnvironmentVariables)
	d.Set("resources", stage.Resources)

	//d.Set("approval",stage.Approval)
	//d.Set("jobs",stage.Jobs)

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


	return nil
}

func retrieveStage(d *schema.ResourceData, meta interface{}) (*gocd.Stage, error) {
	var stageName string
	if name, hasName := d.GetOk("name"); hasName {
		stageName = name.(string)
	} else {
		return nil, errors.New("Missing `name`")
	}

	var stages []*gocd.Stage
	var existing gocd.StageContainer
	var idFormat string

	name, pType, err := pipelineNameType(d)
	if err != nil {
		return nil, err
	}

	client := meta.(*gocd.Client)
	ctx := context.Background()

	if pType == STAGE_TYPE_PIPELINE_TEMPLATE {
		existing, _, err = client.PipelineTemplates.Get(ctx, name)
		idFormat = "template/%s/%s"
	} else if pType == STAGE_TYPE_PIPELINE {
		existing, _, err = client.PipelineConfigs.Get(ctx, name)
		idFormat = "pipeline/%s/%s"
	}
	if err != nil {
		return nil, err
	}

	d.SetId(fmt.Sprintf(idFormat, existing.GetName(), stageName))
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
