package gocd

import (
	"github.com/hashicorp/terraform/helper/schema"
	"encoding/json"
	"strconv"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/drewsonne/go-gocd/gocd"
)

func dataSourceGocdTaskTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGocdJobTemplateRead,
		Schema: map[string]*schema.Schema{
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGocdTaskTemplateRead(d *schema.ResourceData, meta interface{}) error {

	task := gocd.Task{}
	task.Type = d.Get("type").(string)

	jsonDoc, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		// should never happen if the above code is correct
		return err
	}
	jsonString := string(jsonDoc)
	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(hashcode.String(jsonString)))

	return nil
}
