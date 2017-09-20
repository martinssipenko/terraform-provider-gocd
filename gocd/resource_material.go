package gocd

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePipelineMaterial() *schema.Resource {
	return &schema.Resource{
		Create:   resourcePipelineMaterialCreate,
		Read:     resourcePipelineMaterialRead,
		Update:   resourcePipelineMaterialUpdate,
		Delete:   resourcePipelineMaterialDelete,
		Exists:   resourcePipelineMaterialExists,
		Importer: resourcePipelineMaterialStateImport(),
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"attributes": materialsAttributeSchema(),
		},
	}
}

func resourcePipelineMaterialCreate(d *schema.ResourceData, meta interface{}) (err error) {
	return nil

}

func resourcePipelineMaterialRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePipelineMaterialUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	return nil
}

func resourcePipelineMaterialDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePipelineMaterialExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	return true, nil
}

func resourcePipelineMaterialStateImport() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			d.Set("name", d.Id())
			return []*schema.ResourceData{d}, nil
		},
	}
}
