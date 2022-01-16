locals {
  ecs_cluster_name = "microservices-${terraform.workspace}"

  balance_table  = "balance-${terraform.workspace}"
  exchange_table = "currency-exchange-${terraform.workspace}"

  order_creation_queue    = "order-creation-${terraform.workspace}.fifo"
  order_settlement_queue  = "order-settlement-${terraform.workspace}.fifo"
  movement_creation_queue = "movement-creation-${terraform.workspace}.fifo"

  server_cpu         = 512
  server_soft_memory = 512
  server_hard_memory = 1024
  server_name        = "${var.project}-rest-server-${terraform.workspace}"
  server_port        = 3000
  worker_cpu         = 512
  worker_soft_memory = 256
  worker_hard_memory = 1024
  worker_name        = "${var.project}-worker-${terraform.workspace}"
  account_id         = data.aws_caller_identity.current.account_id

  vpc_id     = data.aws_vpc.default_vpc.id
  subnet_ids = tolist(data.aws_subnet_ids.subnets_from_default_vpc.ids)

  sqs_endpoint = "https://sqs.${data.aws_region.current.name}.amazonaws.com/${data.aws_caller_identity.current.account_id}"

  alarm_arn_prefix = "arn:aws:cloudwatch:${var.aws_region}:${local.account_id}:alarm:${var.project}"

  tags = {
    project = var.project
    env     = terraform.workspace
  }
}
