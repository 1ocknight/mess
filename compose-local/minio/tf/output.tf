output "test" {
  value = minio_iam_user.avatar-backend.id
  sensitive = true
}

output "secret" {
  value = minio_iam_user.avatar-backend.secret
  sensitive = true
}