resource "minio_s3_bucket_notification" "avatar-event" {
  bucket = minio_s3_bucket.avatar-bucket.bucket

  queue {
    id = var.kafka_topic_avatar
    queue_arn = "arn:minio:sqs::avatar-events:kafka"
    events = ["s3:ObjectCreated:*", "s3:ObjectRemoved:*"]
  }
}