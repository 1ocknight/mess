resource "keycloak_openid_client" "client-e2e" {
  realm_id  = keycloak_realm.realm-e2e.id
  client_id = "e2e"
  client_secret = "e2e"

  name      = "from e2e tests"
  enabled   = true

  access_type = "CONFIDENTIAL"
  standard_flow_enabled = true
  direct_access_grants_enabled = true

  valid_redirect_uris = [
    "*"
  ]
}

resource "keycloak_openid_client" "client-main" {
  realm_id  = keycloak_realm.realm-main.id
  client_id = "main"
  client_secret = "main"

  name      = "from e2e tests"
  enabled   = true

  access_type = "CONFIDENTIAL"
  standard_flow_enabled = true
  direct_access_grants_enabled = true

  valid_redirect_uris = [
    "*"
  ]
}
