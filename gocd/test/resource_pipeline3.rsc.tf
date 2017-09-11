## START pipeline.test-pipeline0
# CMD terraform import gocd_pipeline.test-pipeline0 "test-pipeline0"
resource "gocd_pipeline" "test-pipeline0" {
  name = "test-pipeline0"
  group = "defaultGroup"
  label_template  = "$${COUNT}"
  materials = [
    {
      type = "git"
      attributes {
        url = "https://github.com/drewsonne/terraform-provider-gocd.git"
        branch = "master"
        auto_update = true
      }
    },
  ]
}

# CMD terraform import gocd_pipeline_stage.test "test"
resource "gocd_pipeline_stage" "test" {
  name = "test"
  pipeline_template = "test-pipeline0"
  fetch_materials = true
  jobs = [
    "${data.gocd_job_definition.test.json}"
  ]
}
data "gocd_job_definition" "test" {
  name = "test"
  tasks = [
    "${data.gocd_task_definition.test-pipeline0_test_test_0.json}",
  ]
}
data "gocd_task_definition" "test-pipeline0_test_test_0" {
  type = "exec"
  run_if = ["success"]
  arguments = [
    "test"]
}

## END