resource "keycloak_realm" "realm-e2e" {
  realm   = "e2e-tf"
  enabled = true
}
