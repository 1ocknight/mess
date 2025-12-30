resource "minio_iam_user" "avatar-backend" {
  name = "avatar-backend"
  secret = "avatar-backend"
}

resource "minio_iam_user_policy_attachment" "avatar-attach" {
  user_name   = minio_iam_user.avatar-backend.name
  policy_name = "readwrite"
  depends_on = [minio_iam_user.avatar-backend]
}

