package gocd

import (
	r "github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func testEnvironmentAssociation(t *testing.T) {
	t.Run("Import", testResourceEnvironmentAssociationImportBasic)
	t.Run("Basic", testResourceEnvironmentAssociationBasic)
}

func testResourceEnvironmentAssociationBasic(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdEnvironmentAssociationDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_environment_association.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceExists("gocd_environment.test-environment"),
					testCheckResourceName("gocd_environment.test-environment",
						"test-environment"),
					testCheckResourceExists("gocd_pipeline.test-pipeline"),
					testCheckResourceName("gocd_pipeline.test-pipeline",
						"test-pipeline"),
					testCheckResourceExists("gocd_environment_association.test-environment-association"),
					testCheckResourceName("gocd_environment_association.test-environment-association",
						"test-environment-association"),
				),
			},
		},
	})
}
