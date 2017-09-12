package gocd

import (
	"fmt"
	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func testCheckResourceExists(resource string) r.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule()
		rs, ok := r.Resources[resource]
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
		r := s.RootModule()
		if rs := r.Resources[resource]; rs.Primary.ID != id {
			return fmt.Errorf("Expected id '%s', got '%s", id, rs.Primary.ID)
		}

		return nil
	}
}
