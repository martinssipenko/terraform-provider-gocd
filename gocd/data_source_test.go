package gocd

import "testing"

func TestDataSource(t *testing.T) {
	t.Run("JobDefinition", testDataSourceJobDefinition)
	t.Run("StageDefinition", testDataSourceGocdStageTemplateRead)
	t.Run("TaskDefinition", testDataSourceTaskDefinition)
}
