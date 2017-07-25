package gocdprovider

import (
	"github.com/drewsonne/gocdsdk"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// I don't know how to get the test case for InternalValidation and the new plugin version.
func Provider() terraform.ResourceProvider {
	return SchemaProvider()
}

func SchemaProvider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"baseurl": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["gocd_baseurl"],
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["username"],
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["password"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gocd_pipeline_template": resourcePipelineTemplate(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gocd_stage_template_definition": dataSourceGocdStageTemplate(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"gocd_baseurl": "URL for the GoCD Server",
		"username":     "User to interact with the GoCD API with.",
		"password":     "Password for User for GoCD API interaction.",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	baseUrl := d.Get("baseurl").(string)
	auth := gocd.Auth{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}
	return gocd.NewClient(baseUrl, &auth, nil), nil
}
