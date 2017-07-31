package gocd

import (
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
	"encoding/json"
	"github.com/drewsonne/go-gocd/gocd"
	"strconv"
	"github.com/hashicorp/terraform/helper/hashcode"
	"fmt"
	"regexp"
)

var dataSourceAwsIamPolicyDocumentVarReplacer = strings.NewReplacer("&{", "${")

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
			"job": {
				Type:     schema.TypeList,
				MinItems: 1,
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

func validateStageAuthorization(v interface{}, k string) (ws []string, errors []error) {
	val := v.(string)
	if !regexp.MustCompile("^[\\w _]+$").MatchString(val) {
		errors = append(errors, fmt.Errorf("%q must contain only alphanumeric caracters and spaces", k))
	}

	return
}

func dataSourceGocdStageTemplateRead(d *schema.ResourceData, meta interface{}) error {
	doc := gocd.Stage{
		Name:     d.Get("name").(string),
		Approval: &gocd.Approval{},
	}

	if manualApproval, hasManualApproval := d.Get("manual_approval").(bool); hasManualApproval && manualApproval {
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

	//var cfgJobs = d.Get("jobs").([]interface{})
	//jobs := make([]gocd.Job, len(cfgJobs))
	//doc.Jobs = jobs
	//
	//for i, jobI := range cfgJobs {
	//	cfgJob := jobI.(map[string]interface{})
	//
	//	job := gocd.Job{
	//		Name: cfgJob["name"].(string),
	//	}
	//	jobs[i] = job
	//}
	//
	jsonDoc, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		// should never happen if the above code is correct
		return err
	}
	jsonString := string(jsonDoc)
	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(hashcode.String(jsonString)))

	return nil
}
