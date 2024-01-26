package pkg

import (
	"context"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/sjxqqq/starknet-go/utils"
)

func GetAccountPrecomputedAddress(provider *rpc.Provider, acnt *account.Account, PubKey string) (*felt.Felt, error) {
	classHash, err := utils.HexToFelt(OpenzeppelinAccountDeployedClassHash)
	if err != nil {
		return nil, err
	}

	PubKeyFelt, err := utils.HexToFelt(PubKey)
	if err != nil {
		return nil, err
	}

	precomputedAddress, err := acnt.PrecomputeAddress(&felt.Zero, PubKeyFelt, classHash, []*felt.Felt{PubKeyFelt})
	if err != nil {
		return nil, err
	}
	return precomputedAddress, nil
}

func DeployAccount(provider *rpc.Provider, accountAddress string, privateKey string) error {
	// new account
	acnt, err := GetNewAccount(provider, accountAddress, privateKey)

	PubKey := GetPublickeyFromPrivateKey(privateKey)

	precomputedAddress, err := GetAccountPrecomputedAddress(provider, acnt, PubKey)
	if err != nil {
		return err
	}

	fmt.Println("precomputedAddress: ", precomputedAddress)
	_, err = TransferToService(precomputedAddress.String())
	fmt.Println("====> : TransferToService")
	if err != nil {
		fmt.Println(err.Error(), "<<<<<")
		// panic(err)
		return err
	}

	classHash, err := utils.HexToFelt(OpenzeppelinAccountDeployedClassHash)
	if err != nil {
		return err
	}
	PubKeyFelt, err := utils.HexToFelt(PubKey)
	if err != nil {
		return err
	}

	// Create transaction data
	tx := rpc.DeployAccountTxn{
		Nonce:               &felt.Zero, // Contract accounts start with nonce zero.
		MaxFee:              new(felt.Felt).SetUint64(4724395326064),
		Type:                rpc.TransactionType_DeployAccount,
		Version:             rpc.TransactionV1,
		Signature:           []*felt.Felt{},
		ClassHash:           classHash,
		ContractAddressSalt: PubKeyFelt,
		ConstructorCalldata: []*felt.Felt{PubKeyFelt},
	}

	err = DeployAccountTranscation(acnt, tx, precomputedAddress)
	if err != nil {
		return err
	}

	return nil
}

func DeployAccountTranscation(acnt *account.Account, tx rpc.DeployAccountTxn, precomputedAddress *felt.Felt) error {
	// Sign the transaction
	err := acnt.SignDeployAccountTransaction(context.Background(), &tx, precomputedAddress)
	if err != nil {
		return err
	}

	// Send transaction to the network
	resp, err := acnt.AddDeployAccountTransaction(context.Background(), rpc.BroadcastDeployAccountTxn{DeployAccountTxn: tx})
	if err != nil {
		return err
		// panic(fmt.Sprint("Error returned from AddDeployAccountTransaction: ", err))
	}
	fmt.Println("AddDeployAccountTransaction response:", resp)

	fmt.Println("TransactionHash: ", resp.TransactionHash, "\tContractAddress:", resp.ContractAddress)
	return nil
}
