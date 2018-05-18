package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	r "github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func testEnvironment(t *testing.T) {
	t.Run("Import", testResourceEnvironmentImportBasic)
	t.Run("Basic", testResourceEnvironment_basic)
}

func testResourceEnvironment_basic(t *testing.T) {
	rInt := acctest.RandInt()

	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdEnvironmentDestroy,
		Steps: []r.TestStep{
			{
				Config:  testAccResource_basic(rInt),
				Destroy: false,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr(
						"gocd_environment.test-cases",
						"name",
						fmt.Sprintf("my_test_environment_%d", rInt),
					),
				),
			},
		},
	})
}

func testAccResource_basic(rInt int) string {
	return fmt.Sprintf(`
resource "gocd_environment" "test-cases" {
  name = "my_test_environment_%d"
}
`, rInt)
}
