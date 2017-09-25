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
			"pipeline": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attributes": materialsAttributeSchema(),
		},
	}
}

func resourcePipelineMaterialImport(d *schema.ResourceData, meta interface{}) (rd []*schema.ResourceData, err error) {
	var pipeline, mType string
	var stubMaterial *gocd.Material
	if pipeline, stubMaterial = parseMaterialId(d.Id()); stubMaterial == nil {
		return nil, fmt.Errorf("Could not find material '%s'", d.Id())
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

	pipeline := d.Get("pipeline").(string)
	mType := d.Get("type").(string)

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()
	ctx := context.Background()

	existing, _, err = client.PipelineConfigs.Get(ctx, pipeline)
	if err != nil {
		return err
	}

	materials := cleanPlaceHolderMaterial(existing.Materials)

	attr, err := extractPipelineMaterialAttributes(mType, d.Get("attributes"))
	if err != nil {
		return err
	}

	newMaterial := gocd.Material{
		Type:       mType,
		Attributes: *attr,
	}

	existing.Materials = append(materials, newMaterial)

	if _, _, err = client.PipelineConfigs.Update(ctx, pipeline, existing); err != nil {
		return nil
	}

	d.SetId(generateMaterialId(&newMaterial, existing.Name))

	return nil
}

func resourcePipelineMaterialRead(d *schema.ResourceData, meta interface{}) (err error) {
	var material, stubMaterial *gocd.Material
	var pipelineName string

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	if pipelineName, stubMaterial = parseMaterialId(d.Id()); stubMaterial == nil {
		return fmt.Errorf("Could not find material '%s'", d.Id())
	}
	if material, err = retrieveMaterial(pipelineName, stubMaterial, client); err != nil {
		return err
	}

	d.SetId(generateMaterialId(material, pipelineName))

	d.Set("pipeline", pipelineName)
	d.Set("type", material.Type)

	materialRaw := readPipelineMaterial(material)
	d.Set("attribute", materialRaw["attributes"])

	return nil
}

func resourcePipelineMaterialUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	return nil
}

func resourcePipelineMaterialDelete(d *schema.ResourceData, meta interface{}) error {
	var stubMaterial *gocd.Material
	var existing *gocd.Pipeline
	var pipelineName string
	var err error

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()
	ctx := context.Background()

	if pipelineName, stubMaterial = parseMaterialId(d.Id()); stubMaterial == nil {
		return fmt.Errorf("Could not find material '%s'", d.Id())
	}

	existing, _, err = client.PipelineConfigs.Get(ctx, pipelineName)
	if err != nil {
		return err
	}

	newMaterials := []gocd.Material{}
	for _, material := range existing.Materials {
		if !material.Equal(stubMaterial) {
			newMaterials = append(newMaterials, material)
		}
	}
	materials := cleanPlaceHolderMaterial(newMaterials)
	if len(materials) == 0 {
		existing.Materials = []gocd.Material{
			*materialPlaceHolder(),
		}
	} else {
		existing.Materials = materials
	}

	if _, _, err := client.PipelineConfigs.Update(ctx, pipelineName, existing); err != nil {
		return err
	}

	return nil
}

func resourcePipelineMaterialExists(d *schema.ResourceData, meta interface{}) (exists bool, err error) {
	pipeline, stubMaterial := parseMaterialId(d.Id())

	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	if _, err = retrieveMaterial(pipeline, stubMaterial, client); err != nil {
		return false, err
	}

	return true, nil
}

func retrieveMaterial(pipeline string, stubMaterial *gocd.Material, client *gocd.Client) (m *gocd.Material, err error) {
	var existing *gocd.Pipeline
	var isMaterial bool
	ctx := context.Background()
	existing, _, err = client.PipelineConfigs.Get(ctx, pipeline)
	for _, material := range existing.Materials {
		if material.Type == stubMaterial.Type {
			switch stubMaterial.Type {
			case "git":
				isMaterial = material.Equal(stubMaterial)
			default:
				return nil, fmt.Errorf("Unexpected material type '%s'", stubMaterial.Type)
			}
			if isMaterial {
				return &material, nil
			}
		}
	}
	return nil, fmt.Errorf("Could not find material with id: `%s/%s`", pipeline, stubMaterial.Type)
}

func cleanPlaceHolderMaterial(materials []gocd.Material) []gocd.Material {
	cleanMaterials := []gocd.Material{}
	for _, material := range materials {
		if !(material.Type == "git" && material.Attributes.Name == PLACEHOLDER_NAME) {
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

func parseMaterialId(id string) (pipeline string, material *gocd.Material) {
	idParts := strings.Split(id, "/")
	material = &gocd.Material{
		Type:       idParts[1],
		Attributes: gocd.MaterialAttributes{},
	}
	switch idParts[1] {
	case "git":
		gitId := strings.Join(idParts[2:], "/")
		gitParts := strings.Split(gitId, ", ")
		material.Attributes.URL = gitParts[0]
		if len(gitParts) > 1 {
			material.Attributes.Branch = gitParts[1]
		}
	default:
		return "", nil
	}
	return idParts[0], material
}

func generateMaterialId(material *gocd.Material, pipelineName string) (id string) {
	var materialId string
	switch material.Type {
	case "git":
		idParts := []string{material.Attributes.URL}
		if material.Attributes.Branch != "" {
			idParts = append(idParts, material.Attributes.Branch)
		}
		materialId = strings.Join(idParts, ", ")
	}
	return fmt.Sprintf(
		"%s/%s/%s",
		pipelineName,
		material.Type,
		materialId,
	)
}
