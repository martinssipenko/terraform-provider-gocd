package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePipelineStageSetPTypeName(d *schema.ResourceData, pType string, name string) error {
	if pType == STAGE_TYPE_PIPELINE {
		d.Set("pipeline", name)
	} else if pType == STAGE_TYPE_PIPELINE_TEMPLATE {
		d.Set("pipeline_template", name)
	} else {
		return fmt.Errorf("Unexpected pipeline type `%s`", pType)
	}
	return nil
}
