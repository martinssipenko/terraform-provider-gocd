package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestDataSourceTaskDefinition(t *testing.T) {

	testSteps := []resource.TestStep{}
	for i := 0; i <= 5; i++ {
		testSteps = append(
			testSteps,
			testStepComparisonCheck(TestStepJsonComparison{
				Id:           "data.gocd_task_definition.test",
				Config:       testFile(fmt.Sprintf("data_source_task_definition.%d.rsc.tf", i)),
				ExpectedJSON: testFile(fmt.Sprintf("data_source_task_definition.%d.rsp.json", i)),
			}),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGocdProviders,
		Steps:     testSteps,
	})
}
