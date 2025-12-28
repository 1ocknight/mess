resource "keycloak_realm" "realm-e2e" {
  realm   = "e2e-realm"
  enabled = true
}

resource "keycloak_realm" "realm-main" {
  realm   = "main-realm"
  enabled = true
}
