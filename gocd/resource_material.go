package gocd

import (
	"context"
	"fmt"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func resourcePipelineMaterial() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineMaterialCreate,
		Read:   resourcePipelineMaterialRead,
		Update: resourcePipelineMaterialUpdate,
		Delete: resourcePipelineMaterialDelete,
		Exists: resourcePipelineMaterialExists,
		Importer: &schema.ResourceImporter{
			State: resourcePipelineMaterialImport,
		},
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

func resourcePipelineMaterialImport(d *schema.ResourceData, meta interface{}) (rd []*schema.ResourceData, err error) {
	var pipeline, mType string
	if pipeline, mType, _, err = parseGoCDPipelineMaterialId(d.Id()); err != nil {
		return nil, err
	}

	d.Set("pipeline", pipeline)
	d.Set("type", mType)

	if err := resourcePipelineMaterialRead(d, meta); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func resourcePipelineMaterialCreate(d *schema.ResourceData, meta interface{}) (err error) {
	var existing *gocd.Pipeline
	var pipeline, mType string
	if pipeline, mType, _, err = parseGoCDPipelineMaterialId(d.Id()); err != nil {
		return err
	}

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()
	ctx := context.Background()

	existing, _, _ = client.PipelineConfigs.Get(ctx, pipeline)

	materials := cleanPlaceHolderMaterial(existing.Materials)

	attr, err := extractPipelineMaterialAttributes(mType, d.Get("Attributes"))
	if err != nil {
		return err
	}

	newMaterial := gocd.Material{
		Type:       mType,
		Attributes: *attr,
	}

	materials = append(materials, newMaterial)
	existing.Materials = materials

	if _, _, err = client.PipelineConfigs.Update(ctx, pipeline, existing); err != nil {
		return nil
	}

	return nil
}

func resourcePipelineMaterialRead(d *schema.ResourceData, meta interface{}) (err error) {
	var material *gocd.Material
	var pipeline, mType, name string

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	if pipeline, mType, name, err = parseGoCDPipelineMaterialId(d.Id()); err != nil {
		return err
	}
	if material, err = retrieveMaterial(pipeline, mType, name, client); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", pipeline, mType, name))

	d.Set("pipeline", pipeline)
	d.Set("type", material.Type)

	materialRaw := readPipelineMaterial(material)
	d.Set("attribute", materialRaw["attributes"])

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

func retrieveMaterial(pipeline string, mType string, name string, client *gocd.Client) (m *gocd.Material, err error) {
	var existing *gocd.Pipeline
	var isMaterial bool
	ctx := context.Background()
	existing, _, err = client.PipelineConfigs.Get(ctx, pipeline)
	for _, material := range existing.Materials {
		if material.Type == mType {
			switch mType {
			case "git":
				isMaterial = material.Attributes.Name == name
			default:
				return nil, fmt.Errorf("Unexpected material type '%s'", mType)
			}
			if isMaterial {
				return &material, nil
			}
		}
	}
	return nil, fmt.Errorf("Could not find material with id: `%s/%s/%s`", pipeline, mType, name)
}

func cleanPlaceHolderMaterial(materials []gocd.Material) []gocd.Material {
	cleanMaterials := []gocd.Material{}
	for _, material := range materials {
		if material.Type != "git" || material.Attributes.Name != PLACEHOLDER_NAME {
			cleanMaterials = append(cleanMaterials, material)
		}
	}
	return cleanMaterials
}

func materialPlaceHolder() *gocd.Material {
	return &gocd.Material{
		Type: "git",
		Attributes: gocd.MaterialAttributes{
			Name:       PLACEHOLDER_NAME,
			URL:        "git@example.com:repo.git",
			AutoUpdate: false,
		},
	}
}

func parseGoCDPipelineMaterialId(id string) (pipeline string, mType string, name string, err error) {
	idParts := strings.Split(id, "/")
	if len(idParts) == 3 {
		return idParts[0], idParts[1], idParts[2], nil
	}

	return "", "", "", fmt.Errorf("could not parse the provided id `%s`", id)
}
