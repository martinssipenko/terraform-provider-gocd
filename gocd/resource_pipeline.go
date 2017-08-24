package gocd

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/drewsonne/go-gocd/gocd"
	"context"
	"errors"
)

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Create:   resourcePipelineCreate,
		Read:     resourcePipelineRead,
		Update:   resourcePipelineUpdate,
		Delete:   resourcePipelineDelete,
		Exists:   resourcePipelineExists,
		Importer: resourcePipelineStateImport(),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"label_template": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_pipeline_locking": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"stages": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"template"},
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					DiffSuppressFunc: supressJSONDiffs,
				},
			},
			"template": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"stages"},
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePipelineCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePipelineRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePipelineDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePipelineExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	var name string
	if ptname, hasName := d.GetOk("name"); hasName {
		name = ptname.(string)
	} else {
		return false, errors.New("`name` can not be empty")
	}

	p, _, err := meta.(*gocd.Client).Pipelines.Get(context.Background(), name, 0)
	exists := (p.Name == name) && (err == nil)
	return exists, err
}

func resourcePipelineStateImport() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			d.Set("name", d.Id())
			return []*schema.ResourceData{d}, nil
		},
	}
}
