provider "symbiosis" {
  api_key = "YOUR_SYMBIOSIS_API_KEY"
}

provider "kubernetes" {
    host = "https://${symbiosis_cluster.production.endpoint}"

    client_certificate = symbiosis_cluster.production.identity.0.certificate
    client_key = symbiosis_cluster.production.identity.0.private_key
    cluster_ca_certificate = symbiosis_cluster.production.identity.0.ca_certificate
}

resource "symbiosis_cluster" "production" {
  name = "production-cluster"
  region = "germany-1"
}

resource "symbiosis_node_pool" "example" {
  cluster = symbiosis_cluster.production.name

  node_type = "general-int-1"
  quantity = 3
}

resource "symbiosis_team_member" "admins" {
  for_each = toset(["sara@mycorp.com", "john@mycorp.com"])
  email = each.value
  role = "ADMIN"
}

resource "symbiosis_cluster_service_account" "example" {
  cluster_name = symbiosis_cluster.production.name
}
