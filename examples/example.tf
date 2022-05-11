provider "symbiosis" {
  api_key = "YOUR_SYMBIOSIS_API_KEY"
}

resource "symbiosis_cluster" "production" {
  name = "production-cluster"
  region = "germany-1"
}

resource "symbiosis_node_pool" "example" {
  cluster = symbiosis_cluster.production.name
  
  node_type = "general-int-1"
  quantity = 6
}

resource "symbiosis_team_member" "admins" {
  for_each = toset(["sara@mycorp.com", "john@mycorp.com"])
  email = each.value
  role = "ADMIN"
}

