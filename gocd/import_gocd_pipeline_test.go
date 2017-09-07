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
	suffix := randomString(10)
	resourceName := fmt.Sprintf("gocd_pipeline.test-%s", suffix)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGocdPipelineConfig(suffix),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testGocdPipelineDestroy(s *terraform.State) error {

	gocdclient := testGocdProvider.Meta().(*gocd.Client)

	root := s.RootModule()
	for _, rs := range root.Resources {
		if rs.Type != "gocd_pipeline" {
			continue
		}

		_, _, err := gocdclient.PipelineConfigs.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("still exists")
		}
	}

	return nil
}

func testGocdPipelineConfig(suffix string) string {
	return strings.Replace(
		testFile("resource_pipeline.0.rsc.tf"),
		"test-pipeline",
		"test-"+suffix,
		-1,
	)
}
