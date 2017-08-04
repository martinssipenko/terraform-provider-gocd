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
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"resources"},
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tab": {
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
			"artifact": {
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
				Type:     schema.TypeMap,
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
			"tabs": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"artifacts": {
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

func dataSourceGocdJobTemplateRead(d *schema.ResourceData, meta interface{}) error {

	tasks := []gocd.Task{}
	for _, rawTask := range d.Get("tasks").([]interface{}) {
		task := gocd.Task{}
		err := json.Unmarshal([]byte(rawTask.(string)), &task)
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
	}

	j := gocd.Job{
		Name:  d.Get("name").(string),
		Tasks: tasks,
	}

	if ric, ok := d.GetOk("run_instance_count"); ok {
		j.RunInstanceCount = int64(ric.(int))
	}

	if to, ok := d.GetOk("time_out"); ok {
		j.RunInstanceCount = int64(to.(int))
	}

	if rscs := decodeConfigStringList(d.Get("resources").([]interface{})); len(rscs) > 0 {
		j.Resources = rscs
	}

	if tabs := decodeConfigStringList(d.Get("resources").([]interface{})); len(tabs) > 0 {
		j.Resources = tabs
	}

	if a := decodeConfigStringList(d.Get("resources").([]interface{})); len(a) > 0 {
		j.Resources = a
	}

	return definitionDocFinish(d, j)
}
