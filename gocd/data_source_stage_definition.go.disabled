package gocd

import (
	"encoding/json"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func dataSourceGocdStageTemplate() *schema.Resource {

	optionalBoolArg := &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}

	stringArg := &schema.Schema{Type: schema.TypeString}

	return &schema.Resource{
		Read: dataSourceGocdStageTemplateRead,
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
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGocdStageTemplateRead(d *schema.ResourceData, meta interface{}) error {
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
	return definitionDocFinish(d, doc)
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
