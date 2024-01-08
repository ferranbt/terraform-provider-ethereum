package ethereum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func checkTransactionDeployed() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(
			"ethereum_transaction.update", "hash"),
		resource.TestCheckResourceAttrSet(
			"ethereum_transaction.update", "gas_used"),
		resource.TestCheckResourceAttrSet(
			"ethereum_transaction.update", "block_num"),
	)
}

func TestAccTransaction_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ethereum_eoa" "account" {
					mnemonic = "test test test test test test test test test test test junk"
				}

				resource "ethereum_transaction" "update" {
					signer = data.ethereum_eoa.account.signer
					to = "0x74B73aC4158B64004F8379966052b215E2A5fc77"
					value = 100
				}
				`,
				Check: checkTransactionDeployed(),
			},
		},
	})
}

func TestAccTransaction_method_fromArtifact(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ethereum_eoa" "account" {
					mnemonic = "test test test test test test test test test test test junk"
				}

				resource "ethereum_contract_deployment" "deploy" {
					signer = data.ethereum_eoa.account.signer
					artifact = "../testcases/out:Inputs"

					input = [
						"0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5",
						"1",
						"0xcc84c3b12f6ae46a791f06a0297bb2d9e60d1d4e0f7c0aff2f5be06cea9189d4",
						jsonencode({
						  "number" = "1"
						})
					]
				}

				resource "ethereum_transaction" "update" {
					signer = data.ethereum_eoa.account.signer
					to = resource.ethereum_contract_deployment.deploy.contract_address

					artifact = "../testcases/out:Inputs"
					method = "applyFunc"

					input = [
						"0x95222290dd7278aa3ddd389cc1e1d165cc4bafe6",
						"2",
						"0xaa84c3b12f6ae46a791f06a0297bb2d9e60d1d4e0f7c0aff2f5be06cea9189d4",
						jsonencode({
						  "number" = "3"
						})
					]
				}
				`,
				Check: checkTransactionDeployed(),
			},
		},
	})
}

func TestAccTransaction_method_fromSignature(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ethereum_eoa" "account" {
					mnemonic = "test test test test test test test test test test test junk"
				}

				resource "ethereum_contract_deployment" "deploy" {
					signer = data.ethereum_eoa.account.signer
					artifact = "../testcases/out:Inputs"

					input = [
						"0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5",
						"1",
						"0xcc84c3b12f6ae46a791f06a0297bb2d9e60d1d4e0f7c0aff2f5be06cea9189d4",
						jsonencode({
						  "number" = "1"
						})
					]
				}

				resource "ethereum_transaction" "update" {
					signer = data.ethereum_eoa.account.signer
					to = resource.ethereum_contract_deployment.deploy.contract_address

					function = "applyFunc(address,uint256,bytes32,(uint256 number))"

					input = [
						"0x95222290dd7278aa3ddd389cc1e1d165cc4bafe6",
						"2",
						"0xaa84c3b12f6ae46a791f06a0297bb2d9e60d1d4e0f7c0aff2f5be06cea9189d4",
						jsonencode({
						  "number" = "3"
						})
					]
				}
				`,
				Check: checkTransactionDeployed(),
			},
		},
	})
}

func TestAccTransaction_Transfer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ethereum_eoa" "account" {
					mnemonic = "test test test test test test test test test test test junk"
				}
				
				resource "ethereum_eoa" "target" {}

				resource "ethereum_transaction" "update" {
					signer = data.ethereum_eoa.account.signer
					to = resource.ethereum_eoa.target.address
					value = "1 gwei"
				}
				`,
				Check: checkTransactionDeployed(),
			},
		},
	})
}

func TestAccTransaction_method_fromSignature(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ethereum_eoa" "account" {
					mnemonic = "test test test test test test test test test test test junk"
				}
				resource "ethereum_contract_deployment" "deploy" {
					signer = data.ethereum_eoa.account.signer
					artifact = "../testcases/out:Inputs"
					input = [
						"0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5",
						"1",
						"0xcc84c3b12f6ae46a791f06a0297bb2d9e60d1d4e0f7c0aff2f5be06cea9189d4",
						jsonencode({
						  "number" = "1"
						})
					]
				}
				resource "ethereum_transaction" "update" {
					signer = data.ethereum_eoa.account.signer
					to = resource.ethereum_contract_deployment.deploy.contract_address
					function = "applyFunc(address,uint256,bytes32,(uint256 number))"
					input = [
						"0x95222290dd7278aa3ddd389cc1e1d165cc4bafe6",
						"2",
						"0xaa84c3b12f6ae46a791f06a0297bb2d9e60d1d4e0f7c0aff2f5be06cea9189d4",
						jsonencode({
						  "number" = "3"
						})
					]
				}
				`,
				Check: checkTransactionDeployed(),
			},
		},
	})
}
