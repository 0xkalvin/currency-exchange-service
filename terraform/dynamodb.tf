resource "aws_dynamodb_table" "currency_exchange_table" {
  name         = local.exchange_table
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "pk"
  range_key    = "sk"

  attribute {
    name = "pk"
    type = "S"
  }

  attribute {
    name = "sk"
    type = "S"
  }

  attribute {
    name = "customer_id"
    type = "S"
  }

  global_secondary_index {
    name            = "customer_id_index"
    hash_key        = "customer_id"
    range_key       = "sk"
    projection_type = "ALL"
  }

  tags = local.tags
}

resource "aws_dynamodb_table" "balance_table" {
  name         = local.balance_table
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "pk"
  range_key    = "sk"

  attribute {
    name = "pk"
    type = "S"
  }

  attribute {
    name = "sk"
    type = "S"
  }

  attribute {
    name = "owner_id"
    type = "S"
  }

  attribute {
    name = "currency_id"
    type = "S"
  }

  global_secondary_index {
    name            = "owner_id_index"
    hash_key        = "owner_id"
    range_key       = "sk"
    projection_type = "ALL"
  }

  global_secondary_index {
    name            = "currency_id_index"
    hash_key        = "currency_id"
    range_key       = "sk"
    projection_type = "ALL"
  }

  tags = local.tags
}
