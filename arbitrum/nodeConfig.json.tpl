{
  "chain": {
    "info-json": "${chain-info}",
    "name": "My Arbitrum L3 Chain"
  },
  "parent-chain": {
    "connection": {
      "url": "https://sepolia-rollup.arbitrum.io/rpc"
    }
  },
  "http": {
    "addr": "0.0.0.0",
    "port": 8449,
    "vhosts": "*",
    "corsdomain": "*",
    "api": [
      "eth",
      "net",
      "web3",
      "arb",
      "debug"
    ]
  },
  "node": {
    "forwarding-target": "",
    "sequencer": {
      "max-tx-data-size": 85000,
      "enable": true,
      "dangerous": {
        "no-coordinator": true
      },
      "max-block-speed": "250ms"
    },
    "delayed-sequencer": {
      "enable": true
    },
    "batch-poster": {
      "max-size": 90000,
      "enable": true,
      "parent-chain-wallet": {
        "private-key": "${batch-poster}"
      }
    },
    "staker": {
      "enable": true,
      "strategy": "MakeNodes",
      "parent-chain-wallet": {
        "private-key": "${staker}"
      }
    },
    "caching": {
      "archive": true
    }
  }
}