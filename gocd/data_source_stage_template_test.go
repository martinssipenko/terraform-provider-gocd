package gocd

import (
	"github.com/hashicorp/terraform/helper/resource"
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
					testTaskDataSourceStateValue(
						"data.gocd_stage_template_definition.test",
						"json",
						testGoCDStageTemplateExpectedJSON,
					),
				),
			},
		},
	})
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
