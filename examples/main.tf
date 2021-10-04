terraform {
  required_providers {
    plusserver = {
      version = "0.1.0"
      source  = "plusserver.com/plusserver/plusserver"
    }
  }
}

provider "plusserver" {
  client_id     = var.client_id
  client_secret = var.client_secret
  username      = var.username
  password      = var.password
  token_url     = var.token_url
  env           = var.env
}

data "plusserver_domain" "barcom" {
  name = "bar.com"
}

output "domain_found" {
  value = data.plusserver_domain.barcom.domains
}

resource "plusserver_domain_record" "foo" {
  domain_id = data.plusserver_domain.barcom.domain_id
  name      = "foobar"
  content   = "1.1.1.1"
  ttl       = 300
  type      = "A"
}
