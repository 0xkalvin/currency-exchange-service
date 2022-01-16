resource "aws_sqs_queue" "order_creation_queue" {
  delay_seconds              = 0
  message_retention_seconds  = 345600
  max_message_size           = 262144
  name                       = local.order_creation_queue
  receive_wait_time_seconds  = 0
  visibility_timeout_seconds = 30
  fifo_queue                 = true
}

resource "aws_sqs_queue" "order_settlement_queue" {
  delay_seconds              = 0
  message_retention_seconds  = 345600
  max_message_size           = 262144
  name                       = local.order_settlement_queue
  receive_wait_time_seconds  = 0
  visibility_timeout_seconds = 30
  fifo_queue                 = true
}

resource "aws_sqs_queue" "movement_creation_queue" {
  delay_seconds              = 0
  message_retention_seconds  = 345600
  max_message_size           = 262144
  name                       = local.movement_creation_queue
  receive_wait_time_seconds  = 0
  visibility_timeout_seconds = 30
  fifo_queue                 = true
}
