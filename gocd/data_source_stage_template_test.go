package gocd

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestDataSourceStageTemplate(t *testing.T) {

	test_steps := []resource.TestStep{}
	for _, test := range []TestStepJsonComparison{
		{
			Id:           "data.gocd_stage_definition.test",
			Config:       testFile("data_source_stage_template.0.rsc.tf"),
			ExpectedJSON: testFile("data_source_stage_template.0.rsp.json"),
		},
	} {
		test_steps = append(
			test_steps,
			testStepComparisonCheck(test),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGocdProviders,
		Steps:     test_steps,
	})
}
