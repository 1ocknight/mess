terraform {
  required_providers {
    kafka = {
      source  = "Mongey/kafka"
      version = "~> 0.7.0"
    }
  }
}

provider "kafka" {
  bootstrap_servers = ["localhost:29092"]
}

resource "kafka_topic" "example" {
  name               = "example-topic"
  partitions         = 3
  replication_factor = 1

  config = {
    cleanup.policy = "delete"
    retention.ms   = "604800000" 
  }
}
