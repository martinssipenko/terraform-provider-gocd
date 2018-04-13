package gocd

import (
	"github.com/beamly/go-gocd/gocd"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testResourcePipelineStageHelpers(t *testing.T) {
	t.Run("ParseJobsFail", testStageParseJobsFail)
	t.Run("ParseIds", testStageParseIDFail)
	t.Run("ParseManualApproval", testParseManualApproval)
}

func testParseManualApproval(t *testing.T) {
	sch := resourcePipelineStage()

	s := &gocd.Stage{}
	ds := sch.Data(&terraform.InstanceState{
		Attributes: map[string]string{
			"authorization_users": "[\"one\"]",
		},
	})
	err := dataSourceStageParseManuallApproval(ds, s)
	assert.Nil(t, err)
}

func testStageParseIDFail(t *testing.T) {
	_, _, _, err := parseGoCDPipelineStageId("not-valid-id")
	assert.EqualError(t, err, "could not parse the provided id `not-valid-id`")
}

func testStageParseJobsFail(t *testing.T) {
	err := dataSourceStageParseJobs([]string{
		"{)",
	}, nil)
	assert.EqualError(t, err, "invalid character ')' looking for beginning of object key string")
}
