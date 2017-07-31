package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"io/ioutil"
	"os"
	"testing"
)

var testGocdProviders map[string]terraform.ResourceProvider
var testGocdProvider *schema.Provider

func init() {
	testGocdProvider = Provider().(*schema.Provider)
	testGocdProviders = map[string]terraform.ResourceProvider{
		"gocd": testGocdProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GOCD_URL"); v == "" {
		t.Fatal("GOCD_URL must be set for acceptance tests.")
	}
	//if v := os.Getenv("GOCD_USERNAME"); v == "" {
	//	t.Fatal("GOCD_USERNAME must be set for acceptance tests.")
	//}
	//if v := os.Getenv("GOCD_PASSWORD"); v == "" {
	//	t.Fatal("GOCD_PASSWORD must be set for acceptance tests.")
	//}

	//var rcfg map[string]interface{}
	//rcfg = make(map[string]interface{})
	//rcfg["baseurl"] = os.Getenv("GOCD_URL")
	//
	//cfg := terraform.ResourceConfig{}
	//cfg.
	//
	//cfg, _ := config.New
	//err := testGocdProvider.Configure(terraform.NewResourceConfig(cfg))
	err := testGocdProvider.Configure(terraform.NewResourceConfig(nil))
	if err != nil {
		t.Fatal(err)
	}
}

// Loads a test file resource from the 'test' directory.
func testFile(name string) string {
	f, err := ioutil.ReadFile(fmt.Sprintf("test/%s", name))
	if err != nil {
		panic(err)
	}

	return string(f)
}
