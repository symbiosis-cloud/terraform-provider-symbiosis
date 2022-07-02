provider "kubernetes" {
    host = "https://${symbiosis_cluster.production.endpoint}"

    client_certificate = symbiosis_cluster.example.identity.0.certificate
    client_key = symbiosis_cluster.example.identity.0.private_key
    cluster_ca_certificate = symbiosis_cluster.example.identity.0.ca_certificate
}

resource "symbiosis_cluster" "example" {
  name = "my-production-cluster"
  region = "germany-1"
}
