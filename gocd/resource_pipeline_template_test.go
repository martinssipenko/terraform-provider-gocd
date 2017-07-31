package gocd

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestResourcePipelineTemplate_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testGocdProviders,
		//CheckDestroy:test
	})
}

//func test
