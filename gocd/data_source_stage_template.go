package gocdprovider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
	"encoding/json"
	"github.com/hashicorp/terraform/helper/hashcode"
	"strconv"
	"github.com/drewsonne/gocdsdk"
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
			"jobs": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"environment_variables": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"approval": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"authorization": {
						Type:     schema.TypeMap,
						Required: true,
						Elem: map[string]*schema.Schema{
							"users": {
								Type:     schema.TypeList,
								Optional: true,
							},
							"roles": {
								Type:     schema.TypeList,
								Optional: true,
							},
						},
					},
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
	doc := gocd.Stage{}
	doc.Name = d.Get("name").(string)
	approval := d.Get("approval").(map[string]interface{})
	if approval != nil {

		doc.Approval = &gocd.Approval{
			Type:          approval["type"].(string),
			Authorization: gocd.Authorization{},
		}

		auth := approval["authorization"].(map[string]string)
		if users, ok := auth["users"]; ok {
			doc.Approval.Authorization.Users = users.([]string)
		} else {
			doc.Approval.Authorization.Roles = auth["roles"]
		}

	}

	var cfgJobs = d.Get("jobs").([]interface{})
	jobs := make([]gocd.Job, len(cfgJobs))
	doc.Jobs = jobs

	for i, jobI := range cfgJobs {
		cfgJob := jobI.(map[string]interface{})

		job := gocd.Job{
			Name: cfgJob["name"].(string),
		}
		jobs[i] = job
	}

	jsonDoc, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		// should never happen if the above code is correct
		return err
	}
	jsonString := string(jsonDoc)
	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(hashcode.String(d.Get("name").(string))))

	return nil
}
