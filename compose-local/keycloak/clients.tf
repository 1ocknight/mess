resource "keycloak_openid_client" "client-e2e" {
  realm_id  = keycloak_realm.realm-e2e.id
  client_id = "e2e"
  client_secret = "fTQwn8jHrT5uJDbn20b1g1jntRFiPU5w"

  name      = "from e2e tests"
  enabled   = true

  access_type = "CONFIDENTIAL"
  standard_flow_enabled = true

  valid_redirect_uris = [
    "*"
  ]
}
