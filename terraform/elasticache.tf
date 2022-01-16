resource "aws_elasticache_subnet_group" "elasticache_subnet_group" {
  name       = "elasticache-sg-${terraform.workspace}"
  subnet_ids = local.subnet_ids
  tags       = local.tags
}


resource "aws_security_group" "elasticache_security_group" {
  vpc_id      = local.vpc_id
  description = "security group for elasticache ${terraform.workspace}"
  name        = "${var.project}-sg-${terraform.workspace}"

  egress {
    from_port = 0
    to_port   = 0
    protocol  = -1

    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.tags
}

resource "aws_security_group_rule" "allow_elasticache_access_for_task" {
  type        = "ingress"
  from_port   = "6379"
  to_port     = "6379"
  protocol    = "tcp"
  description = "Allow inbound elasticache from ${var.project}}"

  source_security_group_id = aws_security_group.server_task_security_group.id

  security_group_id = aws_security_group.elasticache_security_group.id
}

resource "aws_elasticache_cluster" "redis" {
  cluster_id               = "${var.project}-redis-${terraform.workspace}"
  engine                   = "redis"
  node_type                = "cache.t3.micro"
  num_cache_nodes          = 1
  parameter_group_name     = "default.redis6.x"
  engine_version           = "6.x"
  port                     = 6379
  security_group_ids       = [aws_security_group.elasticache_security_group.id]
  subnet_group_name        = aws_elasticache_subnet_group.elasticache_subnet_group.id
  snapshot_retention_limit = 1

  tags = local.tags
}
