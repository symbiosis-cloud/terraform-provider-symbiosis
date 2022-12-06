data "symbiosis_cluster" "cluster" {
    name = "my-cluster"
}

provider "kubernetes" {
    host = "https://${data.symbiosis_cluster.cluster.endpoint}"

    client_certificate     = data.symbiosis_cluster.cluster.certificate
    client_key             = data.symbiosis_cluster.cluster.private_key
    cluster_ca_certificate = data.symbiosis_cluster.cluster.ca_certificate
}