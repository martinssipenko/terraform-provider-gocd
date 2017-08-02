package gocd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePipelineTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineTemplateCreate,
		Read:   resourcePipelineTemplateRead,
		Update: resourcePipelineTemplateUpdate,
		Delete: resourcePipelineTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stages": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourcePipelineTemplateCreate(d *schema.ResourceData, meta interface{}) error {

	var name string
	if ptname, hasName := d.GetOk("name"); hasName {
		name = ptname.(string)
	}

	stages := []*gocd.Stage{}
	for _, rawstage := range d.Get("stages").([]interface{}) {
		stage := gocd.Stage{}
		json.Unmarshal([]byte(rawstage.(string)), &stage)
		stages = append(stages, &stage)
	}

	pt, _, err := meta.(*gocd.Client).PipelineTemplates.Create(context.Background(), name, stages)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("/api/admin/templates/%s", pt.Name))

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
