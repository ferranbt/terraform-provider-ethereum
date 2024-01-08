// Find the first transfer between blocks 10 and 20
data "ethereum_filter_transaction" "filter" {
  is_transfer = true

  start_block  = 10
  limit_blocks = 10
}

// Filter by sender/receiver between blocks 0 and 10
data "ethereum_filter_transaction" "filter" {
  from = "0x.."
  to   = "0x.."

  start_block  = 0
  limit_blocks = 10
}
