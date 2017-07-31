package gocd

import (
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePipelineTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineTemplateCreate,
		Read:   resourcePipelineTemplateRead,
		Update: resourcePipelineTemplateUpdate,
		Delete: resourcePipelineTemplateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stages": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     gocd.Stage{},
			},
		},
	}
}

func resourcePipelineTemplateCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePipelineTemplateRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePipelineTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePipelineTemplateDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
