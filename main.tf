provider "gocd" {
  baseurl = "https://goserver.go.beamly.com:8153/go/"
  username = "notifications"
  password = "notifications"
}


data "gocd_stage_template_definition" "test-stage" {
  name = "test-stage"
  jobs = {
    name = "hallo"
  }
  approval = {
    type = "manual"
    authorization = "one"
  }
}

//resource "gocd_pipeline_template" "my-server" {
//  name = "my-test-template"
//  stages = [
//    ""]
//}

output "stage" {
  value = "${data.gocd_stage_template_definition.test-stage.json}"
}