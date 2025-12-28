output "secret-e2e" {
  value     = keycloak_openid_client.client-e2e.client_secret
  sensitive = true
}

output "id-e2e" {
  value     = keycloak_user.user-e2e.id
  sensitive = true
}

output "secret-main" {
  value     = keycloak_openid_client.client-main.client_secret
  sensitive = true
}
