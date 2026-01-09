//buckets
resource "minio_s3_bucket" "profile-bucket" {
  bucket = "avatar"
  acl    = "private"
}

//users
resource "minio_iam_user" "profile-user" {
  name = "profile"
  secret = "profile-secret"
}
resource "minio_iam_user_policy_attachment" "profile-attach" {
  user_name   = minio_iam_user.profile-user.name
  policy_name = "readwrite"
  depends_on = [minio_iam_user.profile-user]
}

//kafka
resource "minio_s3_bucket_notification" "profile-event" {
  bucket = minio_s3_bucket.profile-bucket.bucket

  queue {
    id = var.kafka_topic_profile
    queue_arn = "arn:minio:sqs::${var.kafka_topic_profile}:kafka"
    events = ["s3:ObjectCreated:*"]
  }
}

//vars
variable "minio-address" {
  type = string
}
variable "kafka_topic_profile" {
  type = string
}

//outputs
output "test" {
  value = minio_iam_user.profile-user.id
  sensitive = true
}
output "secret" {
  value = minio_iam_user.profile-user.secret
  sensitive = true
}