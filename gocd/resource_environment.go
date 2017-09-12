package gocd

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/drewsonne/go-gocd/gocd"
	"context"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvironmentCreate,
		Read:   resourceEnvironmentRead,
		Delete: resourceEnvironmentDelete,
		Exists: resourceEnvironmentExists,
		Importer: &schema.ResourceImporter{
			State: resourceEnvironmentImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceEnvironmentExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	name := d.Get("name").(string)
	client := meta.(*gocd.Client)
	client.Lock()
	defer client.Unlock()

	env, _, err := client.Environments.Get(context.Background(), name)
	exists := (env.Name == name) && (err == nil)
	return exists, err
}

func resourceEnvironmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("name", d.Id())
	return []*schema.ResourceData{d}, nil
}
