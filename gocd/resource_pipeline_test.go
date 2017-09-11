package gocd

import (
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testResourcePipeline(t *testing.T) {
	t.Run("Basic", testResourcePipelineBasic)
	t.Run("ImportBasic", testResourcePipelineImportBasic)
	t.Run("ExistsFail", testResourcePipelineExistsFail)
	t.Run("FullStack", testResourcePipelineFullStack)
}

func testResourcePipelineBasic(t *testing.T) {

	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceExists("gocd_pipeline.test-pipeline"),
					testCheckResourceName(
						"gocd_pipeline.test-pipeline", "pipeline0-terraform"),
				),
			},
			{
				Config: testFile("resource_pipeline.1.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceExists("gocd_pipeline.test-pipeline"),
					testCheckResourceName(
						"gocd_pipeline.test-pipeline", "pipeline0-terraform"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testResourcePipelineFullStack(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline.3.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceExists("gocd_pipeline.test-pipeline3"),
					testCheckResourceName(
						"gocd_pipeline.test-pipeline", "test-pipeline3"),
				),
			},
		},
	})

}

func testResourcePipelineExistsFail(t *testing.T) {
	rd := (&schema.Resource{Schema: map[string]*schema.Schema{
		"name": {Type: schema.TypeString, Required: true},
	}}).Data(&terraform.InstanceState{})

	exists, err := resourcePipelineExists(rd, nil)
	assert.False(t, exists)
	assert.EqualError(t, err, "`name` can not be empty")
}
