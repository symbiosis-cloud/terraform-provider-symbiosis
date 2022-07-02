provider "kubernetes" {
    host = "https://${symbiosis_cluster.production.endpoint}"

    client_certificate = symbiosis_cluster.example.certificate
    client_key = symbiosis_cluster.example.private_key
    cluster_ca_certificate = symbiosis_cluster.example.ca_certificate
}

resource "symbiosis_cluster" "example" {
  name = "my-production-cluster"
  region = "germany-1"
}
