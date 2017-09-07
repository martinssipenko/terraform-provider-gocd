resource "gocd_pipeline" "test-pipeline" {
  name = "pipeline0-terraform"
  group = "testing"
  label_template = "build-$${COUNT}"
  template = "${gocd_pipeline_template.test-pipeline.id}"
  materials = [
    {
      type = "git"
      attributes {
        name = "gocd-github"
        url = "git@github.com:gocd/gocd"
        branch = "feature/my-addition"
        destination = "gocd-dir"
        auto_update = true
      }
    }]
}


resource "gocd_pipeline_template" "test-pipeline" {
  name = "template0-terraform"
}

resource "gocd_pipeline_stage" "test-stage" {
  name = "test-stage"
  fetch_materials = false
  clean_working_directory = false
  never_cleanup_artifacts = false
  approval_success = true
  pipeline_template = "${gocd_pipeline_stage.test-stage.name}"
  jobs = [
    <<JOB
    {
      "name": "job1",
      "tasks": [
        {
          "type": "exec",
          "attributes": {
            "run_if": [
              "passed"
            ],
            "command": "terraform"
          }
        }
      ]
    }
JOB
  ]
}