package gocd

import (
	"sort"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
	"github.com/hashicorp/terraform/helper/hashcode"
	"encoding/json"
)

func decodeConfigStringList(lI []interface{}) []string {
	if len(lI) == 1 {
		return []string{lI[0].(string)}
	}
	ret := make([]string, len(lI))
	for i, vI := range lI {
		ret[i] = vI.(string)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(ret)))
	return ret
}

func definitionDocFinish(d *schema.ResourceData, r interface{}) error {
	jsonDoc, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		// should never happen if the above code is correct
		return err
	}
	jsonString := string(jsonDoc)
	d.Set("json", jsonString)
	d.SetId(strconv.Itoa(hashcode.String(jsonString)))

	return nil

}