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
