package gocd

import (
	r "github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func testEnvironment(t *testing.T) {
	t.Run("Import", testResourceEnvironmentImportBasic)
	t.Run("Basic", testResourceEnvironmentBasic)
}

func testResourceEnvironmentBasic(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdEnvironmentDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_environment.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					testCheckResourceExists("gocd_environment.test-environment"),
					testCheckResourceName(
						"gocd_environment.test-environment", "test-environment"),
				),
			},
		},
	})
}
