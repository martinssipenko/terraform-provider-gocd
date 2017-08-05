package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestDataSourceTaskDefinition(t *testing.T) {

	for i := 0; i <= 5; i++ {
		t.Run(
			fmt.Sprintf("gocd_task_definition.%d", i),
			DataSourceTaskDefinition(t,
				fmt.Sprintf("data_source_task_definition.%d.rsc.tf", i),
				fmt.Sprintf("data_source_task_definition.%d.rsp.json", i),
			),
		)
	}
}

func DataSourceTaskDefinition(t *testing.T, configPath string, expectedPath string) func(t *testing.T) {
	return func(t *testing.T) {
		config := testFile(configPath)
		expected := testFile(expectedPath)
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testGocdProviders,
			Steps: []resource.TestStep{testStepComparisonCheck(&TestStepJSONComparison{
				ID:           "data.gocd_task_definition.test",
				Config:       config,
				ExpectedJSON: expected,
			})},
		})
	}

}
