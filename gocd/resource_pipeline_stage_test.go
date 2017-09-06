package gocd

import (
	"context"
	"fmt"
	"github.com/drewsonne/go-gocd/gocd"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testResourceStage(t *testing.T) {
	t.Run("Basic", testResourceStageBasic)
	t.Run("Import", testResourcePipelineStageImportBasic)
	t.Run("PTypeName", testResourcePipelineStagePtypeName)
}

func testResourcePipelineStagePtypeName(t *testing.T) {
	t.Run("Template", testResourcePipelineStagePtypeNameTemplate)
	t.Run("Pipeline", testResourcePipelineStagePtypeNamePipeline)
	t.Run("Fail", testResourcePipelineStagePtypeNameFail)
}

func testResourcePipelineStagePtypeNameFail(t *testing.T) {
	ds := (&schema.Resource{Schema: map[string]*schema.Schema{
		"pipeline":          {Type: schema.TypeString, Optional: true},
		"pipeline_template": {Type: schema.TypeString, Optional: true},
	}}).Data(&terraform.InstanceState{})
	err := resourcePipelineStageSetPTypeName(ds, "unknown-type", "test-pipeline")
	assert.EqualError(t, err, "Unexpected pipeline type `unknown-type`")

	p, pOk := ds.GetOk("pipeline_template")
	assert.False(t, pOk)
	assert.Empty(t, p)

	pt, ptOk := ds.GetOk("pipeline")
	assert.False(t, ptOk)
	assert.Empty(t, pt)
}

func testResourcePipelineStagePtypeNamePipeline(t *testing.T) {
	ds := (&schema.Resource{Schema: map[string]*schema.Schema{
		"pipeline":          {Type: schema.TypeString, Optional: true},
		"pipeline_template": {Type: schema.TypeString, Optional: true},
	}}).Data(&terraform.InstanceState{})
	err := resourcePipelineStageSetPTypeName(ds, STAGE_TYPE_PIPELINE, "test-pipeline")
	assert.Nil(t, err)

	p, pOk := ds.GetOk("pipeline_template")
	assert.False(t, pOk)
	assert.Empty(t, p)

	pt, ptOk := ds.GetOk("pipeline")
	assert.True(t, ptOk)
	assert.Equal(t, pt, "test-pipeline")
}

func testResourcePipelineStagePtypeNameTemplate(t *testing.T) {
	ds := (&schema.Resource{Schema: map[string]*schema.Schema{
		"pipeline":          {Type: schema.TypeString, Optional: true},
		"pipeline_template": {Type: schema.TypeString, Optional: true},
	}}).Data(&terraform.InstanceState{})
	err := resourcePipelineStageSetPTypeName(ds, STAGE_TYPE_PIPELINE_TEMPLATE, "test-pipeline-template")
	assert.Nil(t, err)

	p, pOk := ds.GetOk("pipeline")
	assert.False(t, pOk)
	assert.Empty(t, p)

	pt, ptOk := ds.GetOk("pipeline_template")
	assert.True(t, ptOk)
	assert.Equal(t, pt, "test-pipeline-template")
}

func testResourceStageBasic(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdStageDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline_stage.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckPipelineStageExists("gocd_pipeline_stage.test-stage"),
				),
			},
		},
	})
}

func testCheckPipelineStageExists(resource string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		rcs := s.RootModule().Resources
		rs, ok := rcs[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No pipeline stage name is set")
		}

		return nil
	}
}

func testGocdStageDestroy(s *terraform.State) error {

	client := testGocdProvider.Meta().(*gocd.Client)

	root := s.RootModule()
	for _, rs := range root.Resources {
		if rs.Type != "gocd_pipeline_stage" {
			continue
		}

		_, _, err := client.PipelineTemplates.Get(context.Background(), rs.Primary.ID)
		//stage := pt.GetStage()
		if err == nil {
			return fmt.Errorf("still exists")
		}
	}

	return nil
}
