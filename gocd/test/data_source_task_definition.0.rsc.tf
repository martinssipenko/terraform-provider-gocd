data "gocd_task_definition" "test" {
  type = "exec"
  run_if = [
    "passed"]
  command = "/usr/local/bin/terraform"
  arguments = [
    "-debug",
    "version"]
  working_directory = "/tmp/"
}