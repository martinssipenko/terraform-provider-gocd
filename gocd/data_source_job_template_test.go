package gocd

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
	"github.com/hashicorp/terraform/terraform"
	"fmt"
)

func TestDataSourceJobTemplate(t *testing.T) {
	for _, test := range []struct {
		Config       string
		ExpectedJSON string
	}{
		{
			Config:       testFile("data_source_job_template.0.rsc.tf"),
			ExpectedJSON: testFile("data_source_job_template.0.rsp.json"),
		},
	} {
		resource.Test(t, resource.TestCase{
			Providers: testGocdProviders,
			Steps: []resource.TestStep{
				{
					Config: test.Config,
					Check: resource.ComposeTestCheckFunc(
						testJobTemplateStateValue(
							"data.gocd_job_definition.test",
							"json",
							test.ExpectedJSON,
						),
					),
				},
			},
		})
	}
}

func testJobTemplateStateValue(id string, name string, value string) resource.TestCheckFunc {
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
			return fmt.Errorf("Value for '%s' is:\n%s\nnot:\n%s", name, v, value)
		}

		return nil
	}
}
