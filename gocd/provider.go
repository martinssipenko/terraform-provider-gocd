package gocd

import (
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
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
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOCD_URL",
				}, nil),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["username"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOCD_USERNAME",
				}, nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: descriptions["password"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOCD_PASSWORD",
				}, nil),
			},
			"skip_ssl_check": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: descriptions["skip_ssl_check"],
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GOCD_SKIP_SSL_CHECK",
				}, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gocd_pipeline_template": resourcePipelineTemplate(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gocd_stage_definition": dataSourceGocdStageTemplate(),
			"gocd_job_definition":   dataSourceGocdJobTemplate(),
		},
		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"baseurl":  "URL for the GoCD Server",
		"username": "User to interact with the GoCD API with.",
		"password": "Password for User for GoCD API interaction.",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	baseUrl := d.Get("baseurl").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	skip_ssl_check := d.Get("skip_ssl_check").(bool)

	if baseUrl == "" {
		baseUrl = os.Getenv("GOCD_URL")
	}
	if username == "" {
		username = os.Getenv("GOCD_USERNAME")
	}
	if password == "" {
		password = os.Getenv("GOCD_PASSWORD")
	}

	return gocd.NewClient(baseUrl, &gocd.Auth{
		Username: username,
		Password: password,
	}, nil, skip_ssl_check), nil
}
