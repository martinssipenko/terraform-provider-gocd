package gocd

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestDataSourceTaskDefinition(t *testing.T) {
	test_steps := []resource.TestStep{}
	for _, test := range []TestStepJsonComparison{
		{
			Id:           "data.gocd_task_definition.test",
			Config:       testFile("data_source_task_definition.0.rsc.tf"),
			ExpectedJSON: testFile("data_source_task_definition.0.rsp.json"),
		},
	} {
		test_steps = append(test_steps, resource.TestStep{
			Config: test.Config,
			Check: resource.ComposeTestCheckFunc(
				testTaskDataSourceStateValue(
					test.Id,
					"json",
					test.ExpectedJSON,
				),
			),
		})
	}

	resource.Test(t, resource.TestCase{
		Providers: testGocdProviders,
		Steps:     test_steps,
	})
}

