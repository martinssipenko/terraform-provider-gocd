package gocd

import (
	"context"
	"fmt"
	"github.com/drewsonne/go-gocd/gocd"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strings"
	"testing"
)

func testResourcePipelineMaterialImportBasic(t *testing.T) {
	suffix := randomString(10)
	rscId := "test-" + suffix
	resourceName := "gocd_pipeline_material." + rscId

	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineMaterialDestroy,
		Steps: []r.TestStep{
			{
				Config: testGocdPipelineMaterialConfig(suffix),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "test-pipeline/" + rscId,
			},
		},
	})
}

func testGocdPipelineMaterialDestroy(s *terraform.State) error {

	gocdclient := testGocdProvider.Meta().(*gocd.Client)

	root := s.RootModule()
	for _, rs := range root.Resources {
		if rs.Type != "gocd_pipeline_material" {
			continue
		}

		pName := rs.Primary.Attributes["pipeline"]
		name := rs.Primary.Attributes["name"]
		mType := rs.Primary.Attributes["type"]

		p, _, err := gocdclient.PipelineConfigs.Get(context.Background(), pName)
		for _, material := range p.Materials {
			if material.Type == mType && material.Attributes.Name == name {
				return fmt.Errorf("still exists")
			}
		}
		if err == nil {
			return fmt.Errorf("still exists")
		}
	}

	return nil
}

func testGocdPipelineMaterialConfig(suffix string) string {
	return strings.Replace(
		testFile("resource_pipeline_material.0.rsc.tf"),
		"test-material",
		"test-"+suffix,
		-1,
	)
}
