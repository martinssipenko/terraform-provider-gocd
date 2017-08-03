package gocd

import (
	"encoding/json"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

func dataSourceGocdStageTemplate() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceGocdStageTemplateRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fetch_materials": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"clean_working_directory": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"never_cleanup_artifacts": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"jobs": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"authorization_roles": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"success_approval", "authorization_users"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"environment_variables": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGocdStageTemplateRead(d *schema.ResourceData, meta interface{}) error {
	doc := gocd.Stage{
		Name:     d.Get("name").(string),
		Approval: &gocd.Approval{},
	}

	if manualApproval, ok := d.GetOk("manual_approval"); ok && manualApproval.(bool) {
		doc.Approval.Type = "manual"
		doc.Approval.Authorization = &gocd.Authorization{}
		if users := d.Get("authorization_users").(*schema.Set).List(); len(users) > 0 {
			doc.Approval.Authorization.Users = decodeConfigStringList(users)
		} else if roles := d.Get("authorization_roles").(*schema.Set).List(); len(roles) > 0 {
			doc.Approval.Authorization.Roles = decodeConfigStringList(roles)
		}
	} else if d.Get("success_approval").(bool) {
		doc.Approval.Type = "success"
		doc.Approval.Authorization = nil
	}

	if jobs := decodeConfigStringList(d.Get("jobs").([]interface{})); len(jobs) > 0 {
		for _, rawjob := range jobs {
			job := gocd.Job{}
			json.Unmarshal([]byte(rawjob), &job)
			doc.Jobs = append(doc.Jobs, &job)
		}
	}

	jsonDoc, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	jsonString := string(jsonDoc)
	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(hashcode.String(jsonString)))

	return nil
}
