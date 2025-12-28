terraform {
  required_providers {
    minio = {
      source  = "aminueza/minio"
      version = "~> 2.0"
    }
  }
}

provider "minio" {
  minio_server   = var.minio-address
  minio_user     = "minioadmin"
  minio_password = "minioadmin"
}
