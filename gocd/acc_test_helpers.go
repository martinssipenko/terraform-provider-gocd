package gocd

import (
	"github.com/hashicorp/terraform/terraform"
	"fmt"
	r "github.com/hashicorp/terraform/helper/resource"
)

func testCheckResourceExists(resource string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("`ID` is not set: %s", resource)
		}

		return nil
	}
}

func testCheckResourceName(resource string, id string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		if rs := s.RootModule().Resources[resource]; rs.Primary.ID != id {
			return fmt.Errorf("Expected id '%s', got '%s", id, rs.Primary.ID)
		}

		return nil
	}
}