[
  {
    "chain-id": ${chainId},
    "parent-chain-id": 421614,
    "chain-name": "My Arbitrum L3 Chain",
    "chain-config": {
      "chainId": ${chainId},
      "homesteadBlock": 0,
      "daoForkBlock": null,
      "daoForkSupport": true,
      "eip150Block": 0,
      "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "eip155Block": 0,
      "eip158Block": 0,
      "byzantiumBlock": 0,
      "constantinopleBlock": 0,
      "petersburgBlock": 0,
      "istanbulBlock": 0,
      "muirGlacierBlock": 0,
      "berlinBlock": 0,
      "londonBlock": 0,
      "clique": {
        "period": 0,
        "epoch": 0
      },
      "arbitrum": {
        "EnableArbOS": true,
        "AllowDebugPrecompiles": false,
        "DataAvailabilityCommittee": false,
        "InitialArbOSVersion": 10,
        "InitialChainOwner": "${owner}",
        "GenesisBlockNum": 0
      }
    },
    "rollup": {
      "bridge": "${bridge}",
      "inbox": "${inbox}",
      "sequencer-inbox": "${sequencer-inbox}",
      "rollup": "${rollup}",
      "validator-utils": "${validator-utils}",
      "validator-wallet-creator": "${validator-wallet-creator}",
      "deployed-at": ${deployed-at}
    }
  }
]