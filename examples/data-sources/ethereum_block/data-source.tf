// return the latest block
data "ethereum_block" "block" {
  tag = "latest"
}

// return the block with the given number
data "ethereum_block" "block" {
  number = 18933464
}

// return the block with the given hash
data "ethereum_block" "block" {
  hash = "0xf7579e5ad2ecc1d267f60f14948ee7dc62d9b79a4e719efae41fb56c7b83a908"
}
