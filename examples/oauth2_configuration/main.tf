terraform {
  required_providers {
    cloudamqp = {
      source  = "cloudamqp/cloudamqp"
      version = "~>1.0"
    }
  }
}

provider "cloudamqp" {
  apikey  = var.cloudamqp_customer_api_key
  baseurl = var.cloudamqp_baseurl
}

resource "cloudamqp_oauth2_configuration" "oauth2_configuration" {
  instance_id               = var.instance_id
  resource_server_id        = "rabbitmq"
  issuer                    = var.issuer_url # Such as https://my-keycloak-server/realms/my-realm
  preferred_username_claims = ["user_name", "email"]
  scope_aliases = {
    "MyKey" = "Myrole2"
  }
  verify_aud      = false
  oauth_client_id = "rabbitmq-management"
  oauth_scopes    = ["email", "profile", "rabbitmq.tag:management"]

  sleep   = 5
  timeout = 600
}
