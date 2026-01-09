terraform {
  required_providers {
    keycloak = {
      source  = "keycloak/keycloak"
      version = "~> 5.0"
    }
    minio = {
      source  = "aminueza/minio"
      version = "~> 3.12.0"
    }
  }
}

provider "keycloak" {
  url       = var.keycloak_url
  realm     = "master"
  client_id = "admin-cli"
  username  = var.keycloak_user
  password  = var.keycloak_password
}

provider "minio" {
  minio_server   = var.minio-address
  minio_user     = "minioadmin"
  minio_password = "minioadmin"
}