resource "symbiosis_cluster" "example" {
  name = "my-production-cluster"
  region = "germany-1"
}

resource "symbiosis_node_pool" "example" {
  cluster = symbiosis_cluster.example.name

  node_type = "general-1"
  quantity = 6
  name = "example-pool"
}
