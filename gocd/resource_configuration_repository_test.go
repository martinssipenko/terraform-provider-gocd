package gocd

import (
	r "github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func testConfigurationRepository(t *testing.T) {
	t.Run("Import", testResourceConfigurationRepositoryImportBasic)
	t.Run("Basic", testResourceConfigurationRepositoryBasic)
}

func testResourceConfigurationRepositoryBasic(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdConfigurationRepositoryDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_configuration_repository.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr(
						"gocd_configuration_repository.test-id",
						"id",
						"test-id",
					),
					r.TestCheckResourceAttr(
						"gocd_configuration_repository.test-plugin_id",
						"name",
						"test-plugin_id",
					),
				),
			},
		},
	})
}
