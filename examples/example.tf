provider "symbiosis" {
  api_key = var.symbiosis_api_key
}

resource "symbiosis_cluster" "production" {
  name = "production-cluster"
  region = "germany-1"
  
  nodes {
    node_type = "int-general-1"
    quantity = 6
  }
}

resource "symbiosis_team_member" "admins" {
  for_each = toset(["sara@mycorp.com", "john@mycorp.com"])
  email = each.value
  role = "ADMIN"
}

