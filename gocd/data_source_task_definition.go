package gocd

import (
	"encoding/json"
	"github.com/drewsonne/go-gocd/gocd"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
)

func dataSourceGocdTaskDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGocdTaskDefinitionRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"run_if": {
				Type:     schema.TypeList,
				MaxItems: 3,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"command": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"arguments": {
				Type:     schema.TypeList,
				MaxItems: 3,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"working_directory": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGocdTaskDefinitionRead(d *schema.ResourceData, meta interface{}) error {

	task := gocd.Task{
		Type:       d.Get("type").(string),
		Attributes: gocd.TaskAttributes{},
	}

	if run_if := decodeConfigStringList(d.Get("run_if").([]interface{})); len(run_if) > 0 {
		task.Attributes.RunIf = run_if
	}

	task_type := d.Get("type").(string)

	if task_type == "exec" {
		dataSourceGocdTaskBuildExec(&task, d)
	}

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

func dataSourceGocdTaskBuildExec(t *gocd.Task, d *schema.ResourceData) {
	if cmd, hasCmd := d.GetOk("command"); hasCmd {
		t.Attributes.Command = cmd.(string)
	}

	if args := decodeConfigStringList(d.Get("arguments").([]interface{})); len(args) > 0 {
		t.Attributes.Arguments = args
	}

	if wd, hasWd := d.GetOk("working_directory"); hasWd {
		t.Attributes.WorkingDirectory = wd.(string)
	}

}
