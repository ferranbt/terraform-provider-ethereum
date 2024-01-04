data "ethereum_event" "res" {
  hash     = "0x..."
  artifact = "../testcases/out:WithEvents"
  event    = "One"
}
