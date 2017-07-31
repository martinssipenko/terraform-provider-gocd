provider "gocd" {
  baseurl = "https://goserver.go.beamly.com:8153/go/"
  username = "notifications"
  password = "notifications"
  skip_ssl_check = true
}

data "gocd_stage_template_definition" "manual-approval" {
  name = "test-stage"
  jobs = {
    name = "hallo"
  }
  manual_approval = true
  authorization_roles = [
    "one",
    "two"]
}

data "gocd_stage_template_definition" "success-approval" {
  name = "test-stage"
  jobs = {
    name = "hallo"
  }
  success_approval = true
}

//resource "gocd_pipeline_template" "my-server" {
//  name = "my-test-template"
//  stages = [
//    "${data.gocd_stage_template_definition.test-stage.json}"]
//}
//




output "manual-approval" {
  value = "${data.gocd_stage_template_definition.manual-approval.json}"
}

output "success-approval" {
  value = "${data.gocd_stage_template_definition.success-approval.json}"
}