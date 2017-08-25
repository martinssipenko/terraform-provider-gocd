package gocd

import (
	"context"
	"errors"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
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
			"parameters": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Optional: true,
			},
			"environment_variables": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type: schema.TypeString,
							// ConflictsWith can only be applied to top level configs.
							// A custom validation will need to be used.
							//ConflictsWith: []string{"encrypted_value"},
							Optional: true,
						},
						"encrypted_value": {
							Type: schema.TypeString,
							// ConflictsWith can only be applied to top level configs.
							// A custom validation will need to be used.
							//ConflictsWith: []string{"value"},
							Optional: true,
						},
						"secure": {
							Type:     schema.TypeBool,
							Default:  false,
							Optional: true,
						},
					},
				},
			},
			"materials": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"attributes": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourcePipelineCreate(d *schema.ResourceData, meta interface{}) error {
	var name, group string
	if ptname, hasName := d.GetOk("name"); hasName {
		name = ptname.(string)
	}

	if ptgroup, hasGroup := d.GetOk("group"); hasGroup {
		group = ptgroup.(string)
	}

	p := extractPipeline(d)
	pt, _, err := meta.(*gocd.Client).Pipelines.Create(context.Background(), p, group)
	return readPipeline(d, pt, err)
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

func extractPipeline(d *schema.ResourceData) *gocd.Pipeline {
	return nil
}

func readPipeline(d *schema.ResourceData, p *gocd.PipelineInstance, err error) error {
	return nil
}
