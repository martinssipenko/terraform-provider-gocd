package gocd

import (
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testResourcePipelineMaterial(t *testing.T) {
	t.Run("Basic", testResourceMaterialBasic)
	t.Run("ImportBasic", testResourcePipelineMaterialImportBasic)
	t.Run("ExistsFail", testResourcePipelineMaterialExistsFail)
	//t.Run("FullStack", testResourcePipelineFullStack)
}

func testResourceMaterialBasic(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline_material.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceExists("gocd_pipeline_material.test-material"),
					testCheckResourceName(
						"gocd_pipeline_material.test-pipeline", "material-terraform"),
				),
			},
			{
				Config: testFile("resource_pipeline_material.1.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceExists("gocd_pipeline_material.test-material"),
					testCheckResourceName(
						"gocd_pipeline_material.test-pipeline", "material0-terraform"),
				),
			},
		},
	})
}

func testResourcePipelineMaterialExistsFail(t *testing.T) {
	rd := (&schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Required: true},
	}}).Data(&terraform.InstanceState{})

	exists, err := resourcePipelineExists(rd, nil)
	assert.False(t, exists)
	assert.EqualError(t, err, "`name` can not be empty")
}
