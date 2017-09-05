package gocd

import (
	"context"
	"fmt"
	"github.com/drewsonne/go-gocd/gocd"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func testResourceStage(t *testing.T) {
	t.Run("Basic", testResourceStageBasic)
	t.Run("Import", testResourcePipelineStageImportBasic)
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
