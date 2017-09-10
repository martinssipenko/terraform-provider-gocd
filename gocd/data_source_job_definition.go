package gocd

import (
	"encoding/json"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGocdJobTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGocdJobTemplateRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tasks": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"run_instance_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"environment_variables": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
			"resources": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"elastic_profile_id"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"elastic_profile_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"resources"},
			},
			"tabs": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"artifacts": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"destination": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"properties": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"xpath": {
							Type:     schema.TypeString,
							Required: true,
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

func dataSourceGocdJobTemplateRead(d *schema.ResourceData, meta interface{}) error {

	tasks := []*gocd.Task{}
	for _, rawTask := range d.Get("tasks").([]interface{}) {
		task := gocd.Task{}
		err := json.Unmarshal([]byte(rawTask.(string)), &task)
		if err != nil {
			return err
		}
		tasks = append(tasks, &task)
	}

	j := gocd.Job{
		Name:  d.Get("name").(string),
		Tasks: tasks,
	}

	if ric, ok := d.GetOk("run_instance_count"); ok {
		j.RunInstanceCount = ric.(int)
	}

	if to, ok := d.GetOk("timeout"); ok {
		j.Timeout = to.(int)
	}

	if envVars, ok := d.Get("environment_variables").([]interface{}); ok && len(envVars) > 0 {
		j.EnvironmentVariables = dataSourceGocdJobEnvVarsRead(envVars)
	}

	if props, ok := d.Get("properties").([]interface{}); ok && len(props) > 0 {
		j.Properties = dataSourceGocdJobPropertiesRead(props)
	}

	if resources := d.Get("resources").(*schema.Set).List(); len(resources) > 0 {
		if rscs := decodeConfigStringList(resources); len(rscs) > 0 {
			j.Resources = rscs
		}
	}

	if resources := d.Get("tabs").(*schema.Set).List(); len(resources) > 0 {
		if rscs := decodeConfigStringList(resources); len(rscs) > 0 {
			j.Tabs = rscs
		}
	}

	if resources := d.Get("artifacts").(*schema.Set).List(); len(resources) > 0 {
		if rscs := decodeConfigStringList(resources); len(rscs) > 0 {
			j.Artifacts = rscs
		}
	}

	return definitionDocFinish(d, j)
}

func dataSourceGocdJobPropertiesRead(rawProps []interface{}) []*gocd.JobProperty {
	props := []*gocd.JobProperty{}
	for _, propRaw := range rawProps {
		propStruct := &gocd.JobProperty{}
		prop := propRaw.(map[string]interface{})

		if name, ok := prop["name"]; ok {
			propStruct.Name = name.(string)
		}

		if name, ok := prop["source"]; ok {
			propStruct.Source = name.(string)
		}

		if name, ok := prop["xpath"]; ok {
			propStruct.XPath = name.(string)
		}
		props = append(props, propStruct)
	}
	return props
}

func dataSourceGocdJobEnvVarsRead(rawEnvVars []interface{}) []*gocd.EnvironmentVariable {
	envVars := []*gocd.EnvironmentVariable{}
	for _, envVarRaw := range rawEnvVars {
		envVarStruct := &gocd.EnvironmentVariable{}
		envVar := envVarRaw.(map[string]interface{})

		if name, ok := envVar["name"]; ok {
			envVarStruct.Name = name.(string)
		}

		if val, ok := envVar["value"]; ok {
			envVarStruct.Value = val.(string)
		}

		if encrypted, ok := envVar["encrypted_value"]; ok {
			envVarStruct.EncryptedValue = encrypted.(string)
		}

		if secure, ok := envVar["secure"]; ok {
			envVarStruct.Secure = secure.(string) == "1"
		}

		envVars = append(envVars, envVarStruct)
	}

	return envVars

}
