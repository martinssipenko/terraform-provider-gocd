package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"io/ioutil"
	"os"
	"testing"
)

var (
	testGocdProviders map[string]terraform.ResourceProvider
	testGocdProvider  *schema.Provider
)

type TestStepJSONComparison struct {
	ID           string
	Config       string
	ExpectedJSON string
}

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

	if u := os.Getenv("GOCD_URL"); u == "" {
		t.Fatal("GOCD_URL must be set for acceptance tests.")
	}

	if s := os.Getenv("GOCD_SKIP_SSL_CHECK"); s == "" {
		t.Fatal("GOCD_SKIP_SSL_CHECK must be set for acceptance tests.")
	}

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

func testTaskDataSourceStateValue(id string, name string, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		root := s.RootModule()
		rs, ok := root.Resources[id]
		if !ok {
			return fmt.Errorf("Not found: %s", id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		v := rs.Primary.Attributes[name]
		if v != value {
			return fmt.Errorf("Value for '%s' is:\n%s\nnot:\n%s", name, v, value)
		}

		return nil
	}
}

func testStepComparisonCheck(test TestStepJsonComparison) resource.TestStep {
	return resource.TestStep{
		Config: test.Config,
		Check: resource.ComposeTestCheckFunc(
			testTaskDataSourceStateValue(
				test.Id,
				"json",
				test.ExpectedJSON,
			),
		),
	}
}
