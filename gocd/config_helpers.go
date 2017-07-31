package gocd

import "sort"

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
