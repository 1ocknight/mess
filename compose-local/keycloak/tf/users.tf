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

resource "keycloak_user" "user-main" {
  realm_id = keycloak_realm.realm-main.id

  username = "main"
  enabled  = true

  first_name = "main"
  last_name  = "main"
  email      = "main@main.main"

  initial_password {
    value     = "main"
    temporary = false
  }
}
