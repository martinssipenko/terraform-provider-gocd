package gocd

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestDataSourceJobTemplate(t *testing.T) {
	for _, test := range []TestStepJsonComparison{
		{
			Id:           "data.gocd_job_definition.test",
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
						testTaskDataSourceStateValue(
							test.Id,
							"json",
							test.ExpectedJSON,
						),
					),
				},
			},
		})
	}
}
