provider "gocd" {
  baseurl = "https://***REMOVED***:8153/go/"
  username = "***REMOVED***"
  password = "***REMOVED***"
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