package gocd

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
	"fmt"
)

func TestDataSourceJobDefinition(t *testing.T) {
	test_steps := []resource.TestStep{}
	for i := 0; i <= 1; i++ {
		test_steps = append(
			test_steps,
			testStepComparisonCheck(TestStepJsonComparison{
				Id:           "data.gocd_job_definition.test",
				Config:       testFile(fmt.Sprintf("data_source_job_template.%d.rsc.tf", i)),
				ExpectedJSON: testFile(fmt.Sprintf("data_source_job_template.%d.rsp.json", i)),
			}),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGocdProviders,
		Steps:     test_steps,
	})
}
