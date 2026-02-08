//realms
resource "keycloak_realm" "realm-main" {
  realm   = "main"
  enabled = true
}

//clients
resource "keycloak_openid_client" "client-main" {
  realm_id  = keycloak_realm.realm-main.id
  client_id = "main"
  client_secret = "main"

  name      = "from e2e tests"
  enabled   = true

  access_type = "CONFIDENTIAL"
  standard_flow_enabled = true
  direct_access_grants_enabled = true
  service_accounts_enabled = true

  valid_redirect_uris = [
    "*"
  ]
}

// Service account client для проверки существования пользователей
resource "keycloak_openid_client" "user-checker-service" {
  realm_id  = keycloak_realm.realm-main.id
  client_id = "user-checker-service"
  client_secret = "user-checker-secret"
  name      = "User Checker Service"
  enabled   = true
  access_type = "CONFIDENTIAL"
  service_accounts_enabled = true
  standard_flow_enabled = false
  direct_access_grants_enabled = false
}

data "keycloak_openid_client" "realm_management" {
  realm_id  = keycloak_realm.realm-main.id
  client_id = "realm-management"
}

data "keycloak_role" "view_users" {
  realm_id  = keycloak_realm.realm-main.id
  client_id = data.keycloak_openid_client.realm_management.id
  name      = "view-users"
}

data "keycloak_role" "query_users" {
  realm_id  = keycloak_realm.realm-main.id
  client_id = data.keycloak_openid_client.realm_management.id
  name      = "query-users"
}

// Роли для service account клиента
resource "keycloak_openid_client_service_account_role" "user_checker_view_users" {
  realm_id                = keycloak_realm.realm-main.id
  service_account_user_id = keycloak_openid_client.user-checker-service.service_account_user_id
  client_id               = data.keycloak_openid_client.realm_management.id
  role                    = data.keycloak_role.view_users.name
}

// Назначаем роль query-users service account'у
resource "keycloak_openid_client_service_account_role" "user_checker_query_users" {
  realm_id                = keycloak_realm.realm-main.id
  service_account_user_id = keycloak_openid_client.user-checker-service.service_account_user_id
  client_id               = data.keycloak_openid_client.realm_management.id
  role                    = data.keycloak_role.query_users.name
}

//users
resource "keycloak_user" "user" {
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
resource "keycloak_user" "user-2" {
  realm_id = keycloak_realm.realm-main.id

  username = "test"
  enabled  = true

  first_name = "test"
  last_name  = "test"
  email      = "test@test.test"

  initial_password {
    value     = "test"
    temporary = false
  }
}

//vars
variable "keycloak_url" {
  type = string
}
variable "keycloak_user" {
  type = string
}
variable "keycloak_password" {
  type      = string
  sensitive = true
}
