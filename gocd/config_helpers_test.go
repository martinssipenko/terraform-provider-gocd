package gocd

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigHelper(t *testing.T) {
	t.Run("DecodeConfigStringList/SuccesSingle", decodeConfigStringListSuccessSingle)
	t.Run("DecodeConfigStringList/SuccesMulti", decodeConfigStringListSuccessMulti)
	t.Run("DecodeConfigStringList/FailInt", decodeConfigStringListFailInt)
	t.Run("DefinitionDocFinish/Success", testDefinitionDocFinishSuccess)
	t.Run("DefinitionDocFinish/Fail", testDefinitionDocFinishFail)
}

func testDefinitionDocFinishFail(t *testing.T) {
	err := definitionDocFinish(
		&schema.ResourceData{},
		make(chan int),
	)
	assert.NotNil(t, err)
}

func testDefinitionDocFinishSuccess(t *testing.T) {
	expectedJson := `{
  "one": "hello",
  "two": "world"
}`
	rd := (&schema.Resource{Schema: map[string]*schema.Schema{
		"json": {Type: schema.TypeString, Computed: true},
	}}).Data(&terraform.InstanceState{})
	st := map[string]string{"one": "hello", "two": "world"}
	err := definitionDocFinish(rd, st)

	assert.Nil(t, err)
	assert.Equal(t, expectedJson, rd.Get("json"))
	assert.Equal(t, "3710939758", rd.Id())
}

func decodeConfigStringListFailInt(t *testing.T) {
	n := []int{6, 7, 8}
	i := make([]interface{}, len(n))
	for j := range n {
		i[j] = n[j]
	}
	assert.Panics(t, func() { decodeConfigStringList(i) })
}

func decodeConfigStringListSuccessSingle(t *testing.T) {
	s := []string{"one"}
	i := make([]interface{}, len(s))
	for j := range s {
		i[j] = s[j]
	}
	strs := decodeConfigStringList(i)

	assert.Len(t, strs, 1)
	assert.Equal(t, strs[0], "one")
}

func decodeConfigStringListSuccessMulti(t *testing.T) {
	s := []string{"one", "two"}
	i := make([]interface{}, len(s))
	for j := range s {
		i[j] = s[j]
	}
	strs := decodeConfigStringList(i)

	assert.Len(t, strs, 2)
	assert.Equal(t, strs[0], "one")
	assert.Equal(t, strs[1], "two")
}
