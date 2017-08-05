package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sergi/go-diff/diffmatchpatch"
	"io/ioutil"
	"net/url"
	"os"
	"testing"
)

var (
	testGocdProviders map[string]terraform.ResourceProvider
	testGocdProvider  *schema.Provider
)

type TestStepJSONComparison struct {
	Index        int
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

func testTaskDataSourceStateValue(id string, name string, value string, index int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		root := s.RootModule()
		rs, ok := root.Resources[id]
		if !ok {
			return fmt.Errorf("In '%d'.\nNot found: %s", index, id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("In '%d'.\nNo ID is set", index)
		}

		v := rs.Primary.Attributes[name]
		if v != value {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(v, value, true)
			err := fmt.Errorf("In '%d'.\nValue mismatch for 'json' is:\n%s", value, dmp.DiffPrettyText(diffs))
			//return fmt.Errorf("In '%d'.\nValue mismatch for 'json' is:\n%s", t.Index, )
			return err
		}

		return nil
	}
}

func testStepComparisonCheck(t *TestStepJSONComparison) resource.TestStep {
	return resource.TestStep{
		Config: t.Config,
		Check: func(s *terraform.State) error {
			root := s.RootModule()
			rs, ok := root.Resources[t.ID]
			if !ok {
				return fmt.Errorf("In '%d'.\nNot found: %s", t.Index, t.ID)
			}
			if rs.Primary.ID == "" {
				return fmt.Errorf("In '%d'.\nNo ID is set", t.Index)
			}

			if v := rs.Primary.Attributes["json"]; v != t.ExpectedJSON {
				dmp := diffmatchpatch.New()
				diffs := dmp.PatchToText(
					dmp.PatchMake(
						dmp.DiffMain(v, t.ExpectedJSON, true),
					),
				)
				rawDiff, _ := url.QueryUnescape(diffs)

				err := fmt.Errorf("In '%d'.\nValue mismatch for 'json' is:\n%s", t.Index, rawDiff)
				//return fmt.Errorf("In '%d'.\nValue mismatch for 'json' is:\n%s", t.Index, )
				return err
			}

			return nil
		},
	}
}
