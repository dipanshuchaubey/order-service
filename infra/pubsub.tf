terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }
}

provider "aws" {
  region = "ap-south-1"
}

resource "aws_sqs_queue" "order_consumer_queue" {
  name       = "order-service-queue.fifo"
  fifo_queue = true


  content_based_deduplication = true
  visibility_timeout_seconds  = 60
  delay_seconds               = 0
  message_retention_seconds   = 86400
  receive_wait_time_seconds   = 20
}

resource "aws_sns_topic" "order_publisher_topic" {
  name                        = "order-service-topic.fifo"
  fifo_topic                  = true
  content_based_deduplication = true
}

resource "aws_sns_topic_subscription" "order_consumer_queue_subscription" {
  topic_arn = aws_sns_topic.order_publisher_topic.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.order_consumer_queue.arn
}

data "aws_iam_policy_document" "order_publisher_policy" {
  statement {
    effect = "Allow"
    sid    = "Allow-SNS-SQS"

    principals {
      type        = "Service"
      identifiers = ["sns.amazonaws.com"]
    }

    actions   = ["sqs:SendMessage"]
    resources = [aws_sqs_queue.order_consumer_queue.arn]

    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"
      values   = [aws_sns_topic.order_publisher_topic.arn]
    }
  }
}

resource "aws_sqs_queue_policy" "order_publisher_allow_sqs" {
  queue_url = aws_sqs_queue.order_consumer_queue.id
  policy    = data.aws_iam_policy_document.order_publisher_policy.json
}
