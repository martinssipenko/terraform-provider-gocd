package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestDataSourceStageTemplate(t *testing.T) {
	for i := 0; i <= 0; i++ {
		t.Run(
			fmt.Sprintf("gocd_stage_definition.%d", i),
			DataSourceStageDefinition(t, i,
				fmt.Sprintf("data_source_stage_definition.%d.rsc.tf", i),
				fmt.Sprintf("data_source_stage_definition.%d.rsp.json", i),
			),
		)
	}
}

func DataSourceStageDefinition(t *testing.T, index int, configPath string, expectedPath string) func(t *testing.T) {
	return func(t *testing.T) {
		config := testFile(configPath)
		expected := testFile(expectedPath)
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testGocdProviders,
			Steps: testStepComparisonCheck(&TestStepJSONComparison{
				Index:        index,
				ID:           "data.gocd_stage_definition.test",
				Config:       config,
				ExpectedJSON: expected,
			}),
		})
	}
}
