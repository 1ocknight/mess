resource "keycloak_user" "user-e2e" {
  realm_id = keycloak_realm.realm-e2e.id

  username = "e2e"
  enabled  = true

  first_name = "e2e"
  last_name  = "e2e"
  email      = "e2e@e2e.e2e"

  initial_password {
    value     = "e2e"
    temporary = false
  }
}