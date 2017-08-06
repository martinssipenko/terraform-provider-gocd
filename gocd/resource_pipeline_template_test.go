package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestResourcePipelineTemplate_Basic(t *testing.T) {
	var out string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testFile("resource_pipeline_template.0.rsc.tf"),
				Check: resource.ComposeTestCheckFunc(
					testCheckPipelineTemplateExists("gocd_pipeline_template.test-pipeline", &out),
					testCheckPipelineTemplateName(
						"gocd_pipeline_template.test-pipeline", "template1"),
					testCheckPipelineTemplate1StageCount(
						"gocd_pipeline_template.test-pipeline", "template1"),
				),
				Destroy: false,
			},
			{
				ExpectNonEmptyPlan: true,
				Config:             testFile("resource_pipeline_template.1.rsc.tf"),
				Check: resource.ComposeTestCheckFunc(
					testCheckPipelineTemplateExists("gocd_pipeline_template.test-pipeline", &out),
					testCheckPipelineTemplateName(
						"gocd_pipeline_template.test-pipeline", "template1"),
					testCheckPipelineTemplate2StageCount(
						"gocd_pipeline_template.test-pipeline", "template1"),
				),
			},
		},
	})

}

func testCheckPipelineTemplate1StageCount(resource string, id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[resource].Primary
		if rs.Attributes["stages.#"] != "1" {
			return fmt.Errorf("Expected 1 stage. Found '%s'", rs.Attributes["stages.#"])
		}
		return nil
	}
}

func testCheckPipelineTemplate2StageCount(resource string, id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[resource].Primary
		if rs.Attributes["stages.#"] != "2" {
			return fmt.Errorf("Expected 2 stages. Found '%s'", rs.Attributes["stages.#"])
		}
		return nil
	}
}

func testCheckPipelineTemplateName(resource string, id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[resource]
		if rs.Primary.ID != "template1-terraform" {
			return fmt.Errorf("Expected id 'template1-terraform', got '%s", rs.Primary.ID)
		}

		return nil
	}
}

func testCheckPipelineTemplateExists(resource string, res *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No pipeline template name is set")
		}

		return nil
	}
}
