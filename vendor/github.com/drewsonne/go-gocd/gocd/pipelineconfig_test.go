package gocd

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPipelineConfig(t *testing.T) {
	setup()
	defer teardown()
	t.Run("Create", testPipelineConfigCreate)
	t.Run("Update", testPipelineConfigUpdate)
	t.Run("Delete", testPipelineConfigDelete)
	t.Run("Get", testPipelineConfigGet)
}

func testPipelineConfigGet(t *testing.T) {
	p, resp, err := client.PipelineConfigs.Get(context.Background(), "mock-pipeline")
	assert.Nil(t, p)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "Not Implemented")
}

func testPipelineConfigDelete(t *testing.T) {

	mux.HandleFunc("/api/admin/pipelines/test-pipeline", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "DELETE", "Unexpected HTTP method")
		assert.Equal(t, r.Header.Get("Accept"), apiV4)

		fmt.Fprint(w, `{
  "message": "Pipeline 'test-pipeline' was deleted successfully."
}`)
	})
	message, resp, err := client.PipelineConfigs.Delete(context.Background(), "test-pipeline")
	if err != nil {
		assert.Error(t, err)
	}
	assert.NotNil(t, resp)
	assert.Equal(t, "Pipeline 'test-pipeline' was deleted successfully.", message)
}

func testPipelineConfigCreate(t *testing.T) {
	mux.HandleFunc("/api/admin/pipelines", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST", "Unexpected HTTP method")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(
			t,
			"{\n  \"group\": \"test-group\",\n  \"pipeline\": {\n    \"name\": \"\",\n    \"stages\": null\n  }\n}\n",
			string(b))
		j, _ := ioutil.ReadFile("test/resources/pipelineconfig.0.json")
		fmt.Fprint(w, string(j))
	})

	p := Pipeline{}
	pgs, _, err := client.PipelineConfigs.Create(context.Background(), "test-group", &p)
	if err != nil {
		t.Error(t, err)
	}

	assert.NotNil(t, pgs)
}

func testPipelineConfigUpdate(t *testing.T) {
	mux.HandleFunc("/api/admin/pipelines/test-name", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method, "Unexpected HTTP method")
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(
			t,
			"{\n  \"pipeline\": {\n    \"name\": \"\",\n    \"stages\": null,\n    \"version\": \"test-version\"\n  }\n}\n",
			string(b))
		j, _ := ioutil.ReadFile("test/resources/pipelineconfig.0.json")
		fmt.Fprint(w, string(j))
	})

	p := Pipeline{
		Version: "test-version",
	}
	pcs, _, err := client.PipelineConfigs.Update(context.Background(), "test-name", &p)
	if err != nil {
		t.Error(t, err)
	}

	assert.NotNil(t, pcs)
}
