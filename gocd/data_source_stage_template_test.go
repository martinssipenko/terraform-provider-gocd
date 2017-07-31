package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestDataSourceStageTemplate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGocdProviders,
		Steps: []resource.TestStep{
			{
				Config: testGoCDStageTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					testGoCDStageTemplateStateValue(
						"data.gocd_stage_template_definition.test",
						"json",
						testGoCDStageTemplateExpectedJSON,
					),
				),
			},
		},
	})
}

func testGoCDStageTemplateStateValue(id, name, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		v := rs.Primary.Attributes[name]
		if v != value {
			return fmt.Errorf("Value for %s is $s, not %s", name, v, value)
		}

		return nil
	}
}

var testGoCDStageTemplateConfig = `
data "gocd_stage_template_definition" "test" {
  name = "stage_name"
  jobs = {
    name = "job1"
  }
  manual_approval = true
  authorization_roles = ["one","two"]
}
`

var testGoCDStageTemplateExpectedJSON = `{
	"name": "stage_name",
	"jobs": [{
		"name": "job1"
	}]
}`
