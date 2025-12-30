resource "minio_s3_bucket" "avatar-bucket" {
  bucket = "avatar-test1"
  acl    = "private"
}