resource "aws_cloudwatch_dashboard" "currency_exchange_dashboard" {
  dashboard_name = "${var.project}-${terraform.workspace}-dashboard"

  dashboard_body = jsonencode({
    "widgets" : [
      {
        "type" : "metric",
        "x" : 0,
        "y" : 0,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/ApplicationELB",
              "HTTPCode_Target_5XX_Count",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}",
              "TargetGroup",
              "${aws_lb_target_group.rest_server_target_group.arn_suffix}"
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Server Target Group 5XX Count",
          "yAxis" : {
            "left" : {
              "label" : "Count"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 12,
        "y" : 0,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/ApplicationELB",
              "HTTPCode_Target_4XX_Count",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}",
              "TargetGroup",
              "${aws_lb_target_group.rest_server_target_group.arn_suffix}"
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Server Target Group 4XX Count",
          "yAxis" : {
            "left" : {
              "label" : "Count"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 0,
        "y" : 6,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/ApplicationELB",
              "HTTPCode_Target_2XX_Count",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}",
              "TargetGroup",
              "${aws_lb_target_group.rest_server_target_group.arn_suffix}"
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Server Target Group 2XX Count",
          "yAxis" : {
            "left" : {
              "label" : "Count"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 12,
        "y" : 6,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/ApplicationELB",
              "HTTPCode_Target_2XX_Count",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}",
              "TargetGroup",
              "${aws_lb_target_group.rest_server_target_group.arn_suffix}",
              {
                "id" : "m1",
                "visible" : false
              }
            ],
            [
              "AWS/ApplicationELB",
              "HTTPCode_Target_4XX_Count",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}",
              "TargetGroup",
              "${aws_lb_target_group.rest_server_target_group.arn_suffix}",
              {
                "id" : "m2",
                "visible" : false
              }
            ],
            [
              "AWS/ApplicationELB",
              "HTTPCode_Target_5XX_Count",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}",
              "TargetGroup",
              "${aws_lb_target_group.rest_server_target_group.arn_suffix}",
              {
                "id" : "m3",
                "visible" : false
              }
            ],
            [
              {
                "expression" : "m1 + m2 + m3",
                "id" : "m4",
                "label" : "2xx + 4xx + 5xx",
                "visible" : false
              }
            ],
            [
              {
                "expression" : "m1 / m4 * 100",
                "id" : "m5",
                "label" : "Success Rate",
                "color" : "#ff7f0e"
              }
            ],
          ],
          "view" : "timeSeries",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Server Target Group Success Rate  (2xx / (2xx + 4xx + 5xx)) * 100",
          "yAxis" : {
            "left" : {
              "label" : "Success Rate (%)"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 0,
        "y" : 12,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/ApplicationELB",
              "RequestCount",
              "TargetGroup",
              "${aws_lb_target_group.rest_server_target_group.arn_suffix}",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}",
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Server Target Group Request Count",
          "yAxis" : {
            "left" : {
              "label" : "Count"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 12,
        "y" : 12,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            ["AWS/ApplicationELB",
              "TargetResponseTime",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}",
              "TargetGroup",
            "${aws_lb_target_group.rest_server_target_group.arn_suffix}"],
            ["...",
            { "stat" : "p99.00" }],
            ["...",
            { "stat" : "p95.00" }],
            ["...",
            { "stat" : "p90.00" }],
            ["...",
            { "stat" : "p50.00" }]
          ],
          "view" : "timeSeries",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "stat" : "p95.00",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Server Target Group Response Time",
          "yAxis" : {
            "left" : {
              "label" : "TargetResponseTime (seconds)"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 0,
        "y" : 18,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/ApplicationELB",
              "HTTPCode_ELB_5XX_Count",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}"
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Load Balancer 5XX Count",
          "yAxis" : {
            "left" : {
              "label" : "Count"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 12,
        "y" : 18,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/ApplicationELB",
              "HTTPCode_ELB_4XX_Count",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}"
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Load Balancer 4XX Count",
          "yAxis" : {
            "left" : {
              "label" : "Count"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 0,
        "y" : 24,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/ApplicationELB",
              "ActiveConnectionCount",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}"
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Load Balancer Active Connection Count",
          "yAxis" : {
            "left" : {
              "label" : "Count"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 12,
        "y" : 24,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/ApplicationELB",
              "RequestCount",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}"
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Load Balancer Request Count",
          "yAxis" : {
            "left" : {
              "label" : "Count"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 0,
        "y" : 30,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            ["AWS/ECS",
              "MemoryUtilization",
              "ServiceName",
              "${aws_ecs_service.ecs_service_rest_server.name}",
              "ClusterName",
            "${local.ecs_cluster_name}"]
          ],
          "view" : "timeSeries",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "stat" : "Average",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Server Average ECS MemoryUtilization",
          "yAxis" : {
            "left" : {
              "label" : "Average MemoryUtilization (%)"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 12,
        "y" : 30,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            ["AWS/ECS",
              "MemoryUtilization",
              "ServiceName",
              "${aws_ecs_service.ecs_service_worker.name}",
              "ClusterName",
            "${local.ecs_cluster_name}"]
          ],
          "view" : "timeSeries",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "stat" : "Average",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Worker Average ECS MemoryUtilization",
          "yAxis" : {
            "left" : {
              "label" : "Average MemoryUtilization (%)"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 0,
        "y" : 36,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            ["AWS/ECS",
              "CPUUtilization",
              "ServiceName",
              "${aws_ecs_service.ecs_service_rest_server.name}",
              "ClusterName",
            "${local.ecs_cluster_name}"]
          ],
          "view" : "timeSeries",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "stat" : "Average",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Server Average ECS CPUUtilization",
          "yAxis" : {
            "left" : {
              "label" : "Average CPUUtilization (%)"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 12,
        "y" : 36,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            ["AWS/ECS",
              "CPUUtilization",
              "ServiceName",
              "${aws_ecs_service.ecs_service_worker.name}",
              "ClusterName",
            "${local.ecs_cluster_name}"]
          ],
          "view" : "timeSeries",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "stat" : "Average",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Worker Average ECS CPUUtilization",
          "yAxis" : {
            "left" : {
              "label" : "Average CPUUtilization (%)"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 0,
        "y" : 42,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/ApplicationELB",
              "UnHealthyHostCount",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}",
              "TargetGroup",
              "${aws_lb_target_group.rest_server_target_group.arn_suffix}"
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Maximum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Server Target Group Unhealthy Count",
          "yAxis" : {
            "left" : {
              "label" : "Count"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 12,
        "y" : 42,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/ApplicationELB",
              "HealthyHostCount",
              "LoadBalancer",
              "${aws_lb.load_balancer.arn_suffix}",
              "TargetGroup",
              "${aws_lb_target_group.rest_server_target_group.arn_suffix}"
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Maximum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Server Target Group Healthy Count",
          "yAxis" : {
            "left" : {
              "label" : "Count"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 0,
        "y" : 48,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/SQS",
              "ApproximateAgeOfOldestMessage",
              "QueueName",
              "${local.order_creation_queue}",
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Maximum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Order Creation Queue ApproximateAgeOfOldestMessage",
          "yAxis" : {
            "left" : {
              "label" : "Max"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 12,
        "y" : 48,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/SQS",
              "ApproximateNumberOfMessagesVisible",
              "QueueName",
              "${local.order_creation_queue}",
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Order Creation Queue ApproximateNumberOfMessagesVisible",
          "yAxis" : {
            "left" : {
              "label" : "Count"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 0,
        "y" : 54,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/SQS",
              "ApproximateAgeOfOldestMessage",
              "QueueName",
              "${local.order_settlement_queue}",
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Maximum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Order Settlement Queue ApproximateAgeOfOldestMessage",
          "yAxis" : {
            "left" : {
              "label" : "Max"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 12,
        "y" : 54,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/SQS",
              "ApproximateNumberOfMessagesVisible",
              "QueueName",
              "${local.order_settlement_queue}",
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Order Settlement Queue ApproximateNumberOfMessagesVisible",
          "yAxis" : {
            "left" : {
              "label" : "Count"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 0,
        "y" : 60,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/SQS",
              "ApproximateAgeOfOldestMessage",
              "QueueName",
              "${local.movement_creation_queue}",
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Maximum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Movement Creation Queue ApproximateAgeOfOldestMessage",
          "yAxis" : {
            "left" : {
              "label" : "Max"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x" : 12,
        "y" : 60,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/SQS",
              "ApproximateNumberOfMessagesVisible",
              "QueueName",
              "${local.movement_creation_queue}",
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Movement Creation Queue ApproximateNumberOfMessagesVisible",
          "yAxis" : {
            "left" : {
              "label" : "Count"
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 0,
        "y":  66,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/DynamoDB",
              "SuccessfulRequestLatency",
              "TableName",
              "${local.exchange_table}",
              "Operation",
              "TransactWriteItems",
              { "stat": "Average", "color": "#0073BB", "label": "TransactWriteItems latency" }
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Average",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Exchange Table TransactWriteItems Latency (ms)",
          "yAxis": {
            "left": {
              "showUnits": true
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 12,
        "y":  66,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/DynamoDB",
              "SuccessfulRequestLatency",
              "TableName",
              "${local.balance_table}",
              "Operation",
              "TransactWriteItems",
              { "stat": "Average", "color": "#0073BB", "label": "TransactWriteItems latency" }
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Average",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Balance Table TransactWriteItems Latency (ms)",
          "yAxis": {
            "left": {
              "showUnits": true
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 0,
        "y":  72,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/DynamoDB",
              "ProvisionedWriteCapacityUnits",
              "TableName",
              "${local.exchange_table}",
              { "label": "Provisioned", "color": "#E02020" }
            ],
            [
              ".",
              "ConsumedWriteCapacityUnits",
              ".",
              ".",
              { "stat": "Sum", "id": "m1", "visible": false }
            ],
            [
              {
                "expression": "m1/PERIOD(m1)",
                "label": "Consumed",
                "id": "e1",
                "color": "#0073BB",
                "region": "${var.aws_region}"
                }
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Average",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Exchange Table Write Usage (units/sec)",
          "yAxis": {
            "left": {
              "showUnits": false
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 12,
        "y":  72,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/DynamoDB",
              "ProvisionedWriteCapacityUnits",
              "TableName",
              "${local.balance_table}",
              { "label": "Provisioned", "color": "#E02020" }
            ],
            [
              ".",
              "ConsumedWriteCapacityUnits",
              ".",
              ".",
              { "stat": "Sum", "id": "m1", "visible": false }
            ],
            [
              {
                "expression": "m1/PERIOD(m1)",
                "label": "Consumed",
                "id": "e1",
                "color": "#0073BB",
                "region": "${var.aws_region}"
                }
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Average",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Balance Table Write Usage (units/sec)",
          "yAxis": {
            "left": {
              "showUnits": false
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 0,
        "y":  78,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/DynamoDB",
              "ThrottledRequests",
              "TableName",
              "${local.exchange_table}",
               "Operation",
              "TransactWriteItems",
              { "color": "#0073BB" }
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Exchange Table Write Throttled Requests (Count)",
          "yAxis": {
            "left": {
              "showUnits": true
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 12,
        "y":  78,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/DynamoDB",
              "ThrottledRequests",
              "TableName",
              "${local.balance_table}",
              "Operation",
              "TransactWriteItems",
              { "color": "#0073BB" }
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Balance Table Write Throttled Requests (Count)",
          "yAxis": {
            "left": {
              "showUnits": true
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 0,
        "y":  84,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/DynamoDB",
              "ConditionalCheckFailedRequests",
              "TableName",
              "${local.exchange_table}",
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Exchange Table Conditional Check Failed (Count)",
          "yAxis": {
            "left": {
              "showUnits": false
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 12,
        "y":  84,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/DynamoDB",
              "ConditionalCheckFailedRequests",
              "TableName",
              "${local.balance_table}",
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Balance Table Conditional Check Failed (Count)",
          "yAxis": {
            "left": {
              "showUnits": false
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 0,
        "y":  90,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/DynamoDB",
              "SystemErrors",
              "TableName",
              "${local.exchange_table}",
              "Operation",
              "TransactWriteItems",
              { "color": "#0073BB" }
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Exchange Table TransactWriteItems 5XX Count",
          "yAxis": {
            "left": {
              "showUnits": false
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 12,
        "y":  90,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/DynamoDB",
              "SystemErrors",
              "TableName",
              "${local.balance_table}",
              "Operation",
              "TransactWriteItems",
              { "color": "#0073BB" }
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Balance Table TransactWriteItems 5XX Count",
          "yAxis": {
            "left": {
              "showUnits": false
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 0,
        "y":  96,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/DynamoDB",
              "TransactionConflict",
              "TableName",
              "${local.exchange_table}",
              { "color": "#0073BB" }
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Exchange Table TransactionConflict Count",
          "yAxis": {
            "left": {
              "showUnits": false
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 12,
        "y":  96,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/DynamoDB",
              "TransactionConflict",
              "TableName",
              "${local.balance_table}",
              { "color": "#0073BB" }
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Sum",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Balance Table TransactionConflict Count",
          "yAxis": {
            "left": {
              "showUnits": false
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 0,
        "y":  102,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/ElastiCache",
              "EngineCPUUtilization",
              "CacheClusterId",
              "${var.project}-redis-${terraform.workspace}",
              { "color": "#0073BB" }
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Average",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Redis ${element(aws_elasticache_cluster.redis.cache_nodes[*].id, 0)} EngineCPUUtilization",
          "yAxis": {
            "left": {
              "showUnits": false
            }
          }
        }
      },
      {
        "type" : "metric",
        "x": 12,
        "y":  102,
        "width":  12,
        "height": 6,
        "properties" : {
          "metrics" : [
            [
              "AWS/ElastiCache",
              "DatabaseMemoryUsagePercentage",
              "CacheClusterId",
              "${var.project}-redis-${terraform.workspace}",
              { "color": "#0073BB" }
            ]
          ],
          "view" : "timeSeries",
          "stat" : "Average",
          "stacked" : false,
          "region" : "${var.aws_region}",
          "period" : 60,
          "legend" : {
            "position" : "right"
          },
          "title" : "Redis ${element(aws_elasticache_cluster.redis.cache_nodes[*].id, 0)} DatabaseMemoryUsagePercentage",
          "yAxis": {
            "left": {
              "showUnits": false
            }
          }
        }
      },
      {
        "type" : "log",
        "x" : 0,
        "y" : 108,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "region" : "${var.aws_region}",
          "title" : "Error Logs From Server",
          "query": "SOURCE '/ecs/${local.server_name}' | fields @message | filter @message like 'error' | parse @message '*error_message*:*,*message*:*,' as a, b, why, c, d, what | stats count(*) as total by what, why | sort total desc | limit 10",
          "view": "table"
        }
      },
      {
        "type" : "log",
        "x" : 12,
        "y" : 108,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "region" : "${var.aws_region}",
          "title" : "Error Logs From Worker",
          "query": "SOURCE '/ecs/${local.worker_name}' | fields @message | filter @message like 'error' | parse @message '*message*:*,*error_message*:*,' as a, b, what, c, d, e, why  | stats count(*) as total by what, why | sort total desc | limit 10",
          "view": "table"
        }
      },
      {
        "type" : "log",
        "x" : 0,
        "y" : 114,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "region" : "${var.aws_region}",
          "title" : "Warn Logs From Server",
          "query": "SOURCE '/ecs/${local.server_name}' | fields @message | filter @message like 'warn' | parse @message '*message*:*,' as a, b, what | stats count(*) as total by what | sort total desc | limit 10",
          "view": "table"
        }
      },
      {
        "type" : "log",
        "x" : 12,
        "y" : 114,
        "width" : 12,
        "height" : 6,
        "properties" : {
          "region" : "${var.aws_region}",
          "title" : "Warn Logs From Worker",
          "query": "SOURCE '/ecs/${local.worker_name}' | fields @message | filter @message like 'warn' | parse @message '*message*:*,' as a, b, what | stats count(*) as total by what | sort total desc | limit 10",
          "view": "table"
        }
      },
    ]
  })
}
