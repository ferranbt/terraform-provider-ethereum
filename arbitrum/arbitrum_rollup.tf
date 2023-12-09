terraform {
  required_providers {
    ethereum = {
      source = "hashicorp.com/edu/ethereum"
    }

    docker = {
      source = "hashicorp.com/edu/docker"
    }

    local = {
      source = "hashicorp.com/edu/local"
    }
  }
}

provider "ethereum" {
  host = "https://arbitrum-sepolia.infura.io/v3/d0181f4c9e474c98921e856804aea3f2"
  // host = "http://localhost:8545"
}

provider "ethereum" {
  alias = "l3"
  host  = "http://localhost:8449"
}

variable "chainId" {
  type    = number
  default = 3453943831
}

variable "minL2BaseFee" {
  type    = number
  default = 100000000
}

data "ethereum_eoa" "account" {
  privkey = "6fd0b749794e67a63154bdbf406b18ebb3b372111472402ed5032a2651704a71"
  // mnemonic = "test test test test test test test test test test test junk"
}

data "ethereum_eoa" "staker" {
  // 0xB1f67b6704E342e04D52D8E13A175767f02D3a40
  privkey = "f4cfff61495bec9a3a094d63f8013ec13ef474e882909bdc7628112343ed7abf"
}

data "ethereum_eoa" "batchPoster" {
  // 0x2B3D37f91E5e32cfe33857CBE6D90bfda5FD7C40
  privkey = "77011b1216069743ae4317b03cd061f38a16df64860f3a9422c03463f0658193"
}

resource "ethereum_transaction" "create_rollup" {
  signer = data.ethereum_eoa.account.signer

  artifact = "./artifacts:RollupCreator"
  method   = "createRollup"
  to       = "0x06e341073b2749e0bb9912461351f716decda9b0"

  input = [
    jsonencode({
      "config" : jsonencode({
        "confirmPeriodBlocks" : 150,
        "extraChallengeTimeBlocks" : 0,
        "stakeToken" : "0x0000000000000000000000000000000000000000",
        "baseStake" : 100000000000000000,
        "wasmModuleRoot" : "0x0754e09320c381566cc0449904c377a52bd34a6b9404432e80afd573b67f7b17",
        "owner" : "${data.ethereum_eoa.account.address}",
        "loserStakeEscrow" : "0x0000000000000000000000000000000000000000",
        "chainId" : var.chainId,
        "chainConfig" : templatefile("${path.module}/chain-config.json.tpl", {
          owner   = "${data.ethereum_eoa.account.address}",
          chainId = var.chainId
        }),
        "genesisBlockNum" : 0,
        "sequencerInboxMaxTimeVariation" : jsonencode({
          "delayBlocks" : 5760,
          "futureBlocks" : 48,
          "delaySeconds" : 86400,
          "futureSeconds" : 3600,
        }),
      }),
      "batchPoster" : "${data.ethereum_eoa.batchPoster.address}",
      "validators" : ["0x46314785c30cCcE5BfBC5670e5034007686166a0"],
      "maxDataSize" : 104857,
      "nativeToken" : "0x0000000000000000000000000000000000000000",
      "deployFactoriesToL2" : false,
      "maxFeePerGasForRetryables" : 100000000,
    })
  ]
}

data "ethereum_event" "rollupCreated" {
  artifact = "./artifacts:RollupCreator"
  hash     = ethereum_transaction.create_rollup.hash
  event    = "RollupCreated"
}

data "ethereum_event" "rollupInitialized" {
  artifact = "./artifacts:RollupCore"
  hash     = ethereum_transaction.create_rollup.hash
  event    = "RollupInitialized"
}

/*
resource "ethereum_transaction" "create_bridge" {
	signer = data.ethereum_eoa.account.signer

  	artifact = "./artifacts:BridgeCreator"
  	method = "createTokenBridge"
  	to = "0x1c608642d0944e95957a7ac3a478ec17fa191e9a"

	// TODO: Gas estimation fails?
	gas_limit = 40452528
	// TODO: Why this value transfer?
	value = "37307785000000000"

	input = [
		"${data.ethereum_event.rollupCreated.logs.inboxAddress}",
		"${data.ethereum_eoa.account.address}",
		"30452528", // This is hardcoded
		"100000000"
	]
}
*/

resource "local_file" "foo" {
  content = templatefile("${path.module}/nodeConfig.json.tpl", {
    batch-poster = "${data.ethereum_eoa.batchPoster.signer}",
    staker       = "${data.ethereum_eoa.staker.signer}",
    chain-info = replace(replace(templatefile("${path.module}/chain-info.json.tpl", {
      chainId                  = var.chainId,
      owner                    = "${data.ethereum_eoa.account.address}",
      bridge                   = "${data.ethereum_event.rollupCreated.logs.bridge}",
      inbox                    = "${data.ethereum_event.rollupCreated.logs.inboxAddress}",
      sequencer-inbox          = "${data.ethereum_event.rollupCreated.logs.sequencerInbox}",
      rollup                   = "${data.ethereum_event.rollupInitialized.address}",
      validator-utils          = "${data.ethereum_event.rollupCreated.logs.validatorUtils}",
      validator-wallet-creator = "${data.ethereum_event.rollupCreated.logs.validatorWalletCreator}",
      deployed-at              = "${ethereum_transaction.create_rollup.block_num}"
    }), "\"", "\\\""), "\n", "")
  })
  filename = "${path.module}/data/nodeConfig.json"
}

resource "docker_container" "nitro" {
  name = "nitro"

  // offchainlabs/nitro-node:v2.1.3-e815395
  image = "sha256:e9cd4f661e4896b0c21b8f42122704a0540acb71248b7fc967c107778024cdb3"

  ports {
    external = 8449
    internal = 8449
  }

  mounts {
    target = "/tmp/nodeConfig.json"
    source = abspath(local_file.foo.filename)
    type   = "bind"
  }

  command = [
    "--conf.file", "/tmp/nodeConfig.json"
  ]
}

// Fund staker
resource "ethereum_transaction" "fundStaker" {
  signer = data.ethereum_eoa.account.signer
  to     = data.ethereum_eoa.staker.signer
  value  = "0.2 ether"
}

// Fund batch poster
resource "ethereum_transaction" "fundBatchPoster" {
  signer = data.ethereum_eoa.account.signer
  to     = data.ethereum_eoa.batchPoster.signer
  value  = "0.2 ether"
}

// ---- L3 initialization ----

/*
// TODO: We need to do something to wait until the balance is updated
resource "ethereum_transaction" "fundSignerOnL3" {
  signer   = data.ethereum_eoa.account.signer
  to       = data.ethereum_event.rollupCreated.logs.inboxAddress
  function = "function depositEth() public payable"
  value    = "0.1 ether"
}
*/

resource "ethereum_transaction" "setMinimumL2BaseFee" {
  provider = "ethereum.l3"
  signer   = data.ethereum_eoa.account.signer

  artifact = "./artifacts:ArbOwner"
  method   = "setMinimumL2BaseFee"
  to       = "0x0000000000000000000000000000000000000070"

  input = [
    var.minL2BaseFee
  ]
}

resource "ethereum_transaction" "setNetworkFeeAccount" {
  provider = "ethereum.l3"
  signer   = data.ethereum_eoa.account.signer

  artifact = "./artifacts:ArbOwner"
  method   = "setNetworkFeeAccount"
  to       = "0x0000000000000000000000000000000000000070"

  input = [
    data.ethereum_eoa.account.address
  ]
}

resource "ethereum_transaction" "setInfraFeeAccount" {
  provider = "ethereum.l3"
  signer   = data.ethereum_eoa.account.signer

  artifact = "./artifacts:ArbOwner"
  method   = "setInfraFeeAccount"
  to       = "0x0000000000000000000000000000000000000070"

  input = [
    data.ethereum_eoa.account.address
  ]
}

data "ethereum_call" "getL1BaseFeeEstimate" {
  artifact = "./artifacts:ArbGasInfo"
  method   = "getL1BaseFeeEstimate"
  to       = "0x000000000000000000000000000000000000006c"
}

data "ethereum_getGasPrice" "gas_price" {
}

resource "ethereum_transaction" "setL1PricePerUnit" {
  provider = "ethereum.l3"
  signer   = data.ethereum_eoa.account.signer

  artifact = "./artifacts:ArbOwner"
  method   = "setL1PricePerUnit"
  to       = "0x0000000000000000000000000000000000000070"

  input = [
    tonumber(data.ethereum_call.getL1BaseFeeEstimate.output.0) + data.ethereum_getGasPrice.gas_price.gas_price
  ]
}

// -- transfer ownership --
