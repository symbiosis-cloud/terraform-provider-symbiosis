provider "symbiosis" {
  api_key = "YOUR_SYMBIOSIS_API_KEY"
}

provider "kubernetes" {
    host = "https://${symbiosis_cluster.example.endpoint}"

    client_certificate = symbiosis_cluster.example.certificate
    client_key = symbiosis_cluster.example.private_key
    cluster_ca_certificate = symbiosis_cluster.example.ca_certificate
}

resource "symbiosis_cluster" "example" {
  name = "my-production-cluster"
  region = "germany-1"
}

resource "symbiosis_node_pool" "example" {
  cluster = symbiosis_cluster.example.name

  node_type = "general-1"
  quantity = 3
}
