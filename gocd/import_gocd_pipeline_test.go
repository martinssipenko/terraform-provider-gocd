package gocd

import (
	"context"
	"fmt"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strings"
	"testing"
)

func testResourcePipelineImportBasic(t *testing.T) {
	for _, idx := range []int{4} { //{2,4}
		suffix := randomString(10)
		resourceName := fmt.Sprintf("gocd_pipeline.test-%s", suffix)

		resource.Test(t, resource.TestCase{
			PreCheck:     func() { testAccPreCheck(t) },
			Providers:    testGocdProviders,
			CheckDestroy: testGocdPipelineDestroy,
			Steps: []resource.TestStep{
				{
					Config: testGocdPipelineConfig(suffix, idx),
				},
				{
					ResourceName:      resourceName,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	}
}

func testGocdPipelineDestroy(s *terraform.State) error {

	gocdclient := testGocdProvider.Meta().(*gocd.Client)

	root := s.RootModule()
	for _, rs := range root.Resources {
		if rs.Type != "gocd_pipeline" {
			continue
		}

		if _, _, err := gocdclient.PipelineConfigs.Get(context.Background(), rs.Primary.ID); err == nil {
			return fmt.Errorf("still exists")
		}
	}

	return nil
}

func testGocdPipelineConfig(suffix string, idx int) string {
	return strings.Replace(
		testFile(fmt.Sprintf("resource_pipeline.%d.rsc.tf", idx)),
		"test-pipeline",
		"test-"+suffix,
		-1,
	)
}
