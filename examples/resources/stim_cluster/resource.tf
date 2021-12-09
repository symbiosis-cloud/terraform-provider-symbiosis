resource "stim_cluster" {
  name = "my-staging-cluster"
  region = "eu-germany-1"
  
  nodes {
    node_type = "int-general-1"
    quantity = 6
  }

  nodes {
    node_type = "int-memory-2"
    quantity = 10
  }
}
