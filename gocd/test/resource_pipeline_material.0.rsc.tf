resource "gocd_pipeline" "test-pipeline" {
  name = "test-pipeline"
  group = "test-group"
}

//resource "gocd_pipeline_material" "test-material" {
//  pipeline = "${gocd_pipeline.test-pipeline.name}"
//  type = "git"
//  attributes {
//    name = "test-git-material"
//    url = "https://github.com/gocd/gocd"
//  }
//}