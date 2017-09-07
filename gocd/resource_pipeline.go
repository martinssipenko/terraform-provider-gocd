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
			"group": {
				Type:     schema.TypeString,
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
			"template": {
				Type:     schema.TypeString,
				Optional: true,
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
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"branch": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"destination": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"url": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"auto_update": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"filter": {
										Type:     schema.TypeList,
										Optional: true,
										//Elem: &schema.Schema{
										//	Type: schema.TypeString,
										//},
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ignore": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourcePipelineCreate(d *schema.ResourceData, meta interface{}) error {
	var group string
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

	if p, _, err := meta.(*gocd.Client).PipelineConfigs.Get(context.Background(), name); err != nil {
		return false, err
	} else {
		return (p.Name == name), nil
	}
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
	p := gocd.Pipeline{}

	if template, hasTemplate := d.GetOk("template"); hasTemplate {
		p.Template = template.(string)
	}

	p.Name = d.Get("name").(string)

	rawMaterials := d.Get("materials")
	if materials := rawMaterials.([]interface{}); len(materials) > 0 {
		p.Materials = extractPipelineMaterials(materials)
	}

	return &p
}

func extractPipelineMaterials(rawMaterials []interface{}) []gocd.Material {
	ms := []gocd.Material{}
	for _, rawMaterial := range rawMaterials {
		mat := rawMaterial.(map[string]interface{})
		m := gocd.Material{}

		if mType, ok := mat["type"]; ok {
			m.Type = mType.(string)
		}

		if mAttributes, ok := mat["attributes"]; ok {

			attr := gocd.MaterialAttributes{}

			rawAttr := mAttributes.([]interface{})[0].(map[string]interface{})
			for attrKey, attrValue := range rawAttr {
				switch attrKey {
				case "name":
					attr.Name = attrValue.(string)
				case "url":
					attr.URL = attrValue.(string)
				case "branch":
					attr.Branch = attrValue.(string)
				case "destination":
					attr.Destination = attrValue.(string)
				case "auto_update":
					attr.AutoUpdate = attrValue.(bool)
				case "filter":
					attr.Filter = extractPipelineMaterialFilter(attrValue)
				}
			}

			m.Attributes = attr
		}
		ms = append(ms, m)

	}
	return ms
}

func extractPipelineMaterialFilter(attr interface{}) *gocd.MaterialFilter {
	filterI := attr.([]interface{})[0].(map[string]interface{})
	filters := filterI["ignore"].([]interface{})
	mf := gocd.MaterialFilter{
		Ignore: decodeConfigStringList(filters),
	}
	return &mf
}

func readPipeline(d *schema.ResourceData, p *gocd.Pipeline, err error) error {
	if err != nil {
		return err
	}

	d.SetId(p.Name)
	//d.Set("Version",p.)
	return nil
}
