resource "symbiosis_cluster" {
  name = "my-staging-cluster"
  region = "germany-1"
  
  node_pool {
    node_type = "int-general-1"
    quantity = 6
  }

  node_pool {
    node_type = "int-memory-2"
    quantity = 10
  }
}
