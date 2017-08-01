package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestDataSourceTaskDefinition(t *testing.T) {
	for _, test := range []struct {
		Config       string
		ExpectedJSON string
	}{
		{
			Config:       testFile("data_source_task_definition.0.rsc.tf"),
			ExpectedJSON: testFile("data_source_task_definition.0.rsp.json"),
		},
	} {
		resource.Test(t, resource.TestCase{
			Providers: testGocdProviders,
			Steps: []resource.TestStep{
				{
					Config: test.Config,
					Check: resource.ComposeTestCheckFunc(
						testTaskDefinitionStateValue(
							"data.gocd_task_definition.test-task-exec",
							"json",
							test.ExpectedJSON,
						),
					),
				},
			},
		})
	}
}

func testTaskDefinitionStateValue(id string, name string, value string) resource.TestCheckFunc {
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
