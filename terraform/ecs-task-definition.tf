resource "aws_cloudwatch_log_group" "rest_server_log_group" {
  name              = "/ecs/${local.server_name}"
  retention_in_days = "7"
  tags              = local.tags
}

resource "aws_cloudwatch_log_group" "worker_log_group" {
  name              = "/ecs/${local.worker_name}"
  retention_in_days = "7"
  tags              = local.tags
}

resource "aws_ecs_task_definition" "task_definition_rest_server" {
  family                   = local.server_name
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  execution_role_arn       = aws_iam_role.server_task_iam_role.arn
  task_role_arn            = aws_iam_role.server_task_iam_role.arn
  cpu                      = local.server_cpu
  memory                   = local.server_hard_memory
  container_definitions = jsonencode(
    [
      {
        name : local.server_name,
        image : "${aws_ecr_repository.ecr_repo.repository_url}:rest-server-latest",
        memory: local.server_hard_memory
        memoryReservation: local.server_soft_memory
        environment : [
          {
            name  = "ENV"
            value = "production"
          },
          {
            name  = "PORT"
            value = "3000"
          },
          {
            name  = "DYNAMODB_REGION"
            value = data.aws_region.current.name
          },
          {
            name  = "BALANCE_TABLE"
            value = local.balance_table
          },
          {
            name  = "EXCHANGE_TABLE"
            value = local.exchange_table
          },
          {
            name  = "SQS_ENDPOINT"
            value = local.sqs_endpoint
          },
          {
            name  = "SQS_REGION"
            value = data.aws_region.current.name
          },
          {
            name  = "SQS_CONCURRENCY"
            value = "10"
          },
          {
            name  = "SQS_MAX_RETRIES"
            value = "3"
          },
          {
            name  = "SQS_ORDER_CREATION_QUEUE_URL"
            value = aws_sqs_queue.order_creation_queue.url
          },
          {
            name  = "APP_NAME"
            value = "rest_server"
          },
          {
            name  = "LOG_LEVEL"
            value = "debug"
          }
        ],
        secrets : [],
        essential : true,
        portMappings : [
          {
            containerPort : local.server_port,
            protocol : "tcp"
          }
        ],
        logConfiguration : {
          logDriver : "awslogs",
          options : {
            awslogs-group : "/ecs/${local.server_name}",
            awslogs-region : "${var.aws_region}",
            awslogs-stream-prefix : "ecs"
          }
        }
      }
    ]
  )
  tags = local.tags
}

resource "aws_ecs_task_definition" "task_definition_worker" {
  family                   = local.worker_name
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  execution_role_arn       = aws_iam_role.server_task_iam_role.arn
  task_role_arn            = aws_iam_role.server_task_iam_role.arn
  cpu                      = local.worker_cpu
  memory                   = local.worker_hard_memory
  container_definitions = jsonencode(
    [
      {
        name : local.worker_name,
        image : "${aws_ecr_repository.ecr_repo.repository_url}:worker-latest",
        memory: local.worker_hard_memory
        memoryReservation:  local.worker_soft_memory
        environment : [
          {
            name  = "NODE_ENV"
            value = "production"
          },
          {
            name  = "DYNAMODB_REGION"
            value = data.aws_region.current.name
          },
          {
            name  = "BALANCE_TABLE"
            value = local.balance_table
          },
          {
            name  = "EXCHANGE_TABLE"
            value = local.exchange_table
          },
          {
            name  = "SQS_ENDPOINT"
            value = local.sqs_endpoint
          },
          {
            name  = "SQS_REGION"
            value = data.aws_region.current.name
          },
          {
            name  = "SQS_CONCURRENCY"
            value = "10"
          },
          {
            name  = "SQS_MAX_RETRIES"
            value = "3"
          },
          {
            name  = "SQS_ORDER_CREATION_QUEUE_URL"
            value = aws_sqs_queue.order_creation_queue.url
          },
          {
            name  = "SQS_ORDER_SETTLEMENT_QUEUE_URL"
            value = aws_sqs_queue.order_settlement_queue.url
          },
          {
            name  = "SQS_MOVEMENT_CREATION_QUEUE_URL"
            value = aws_sqs_queue.movement_creation_queue.url
          },
          {
            name  = "SQS_ORDER_CREATION_POINTS"
            value = "200"
          },
          {
            name  = "SQS_ORDER_CREATION_DURATION"
            value = "60"
          },
          {
            name  = "REDIS_HOST"
            value = element(aws_elasticache_cluster.redis.cache_nodes[*].address, 0)
          },
          {
            name  = "REDIS_PORT"
            value = "6379"
          },
          {
            name  = "APP_NAME"
            value = "exchange_worker"
          },
          {
            name  = "LOG_LEVEL"
            value = "debug"
          }
        ],
        secrets : [],
        essential : true,
        logConfiguration : {
          logDriver : "awslogs",
          options : {
            awslogs-group : "/ecs/${local.worker_name}",
            awslogs-region : "${var.aws_region}",
            awslogs-stream-prefix : "ecs"
          }
        }
      }
    ]
  )
  tags = local.tags
}
