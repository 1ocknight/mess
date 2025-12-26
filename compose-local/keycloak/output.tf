output "secret-e2e" {
  value     = keycloak_openid_client.client-e2e.client_secret
  sensitive = true
}
