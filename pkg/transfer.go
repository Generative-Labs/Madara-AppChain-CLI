package pkg

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

func TransferToService(precomputedAddress string) (string, error) {
	targetRpcUrl := "http://127.0.0.1:9944"

	// account_address := ""
	// privateKey := "0x00c1cf1490de1352865301bb8705143f3ef938f97fdf892f1090dcb5ac7bcd1d"

	account_address := "0x3"
	privateKey := "0x00c1cf1490de1352865301bb8705143f3ef938f97fdf892f1090dcb5ac7bcd1d"
	tokenAddress := "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"

	toAccountID := precomputedAddress
	amount := "0x574fbde6000"
	result, err := TransferInvoke(account_address, privateKey, toAccountID, amount, tokenAddress, targetRpcUrl)

	return result, err
}

func TransferInvoke(account_address string, privateKey string, toAccountID string, amount string, tokenAddress string, targetRpcUrl string) (string, error) {
	accountContractVersion := 0

	contractFelt, err := utils.HexToFelt(tokenAddress)
	if err != nil {
		return "", err
	}

	pubkey := GetPublickeyFromPrivateKey(privateKey)

	context := context.Background()
	c, err := ethrpc.DialContext(context, targetRpcUrl)
	if err != nil {
		panic("You need to specify the testnet url")
	}
	clientv02 := rpc.NewProvider(c)

	accountAddressFelt, err := new(felt.Felt).SetString(account_address)
	if err != nil {
		panic("Error casting account_address to felt")
	}

	// Initializing the account memkeyStore
	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		return "", errors.New("invalid Private Key")
	}

	ks.Put(pubkey, fakePrivKeyBI)

	// Set up account
	sAccount, err := account.NewAccount(clientv02, accountAddressFelt, pubkey, ks)
	if err != nil {
		panic(err)
	}

	recipientFelt, err := utils.HexToFelt(toAccountID)
	if err != nil {
		return "", err
	}

	// Starknet token address
	amountFelt, err := utils.HexToFelt(amount)
	if err != nil {
		return "", err
	}

	callData := []*felt.Felt{
		recipientFelt,
		amountFelt,
	}
	// invoke

	// Getting the nonce from the account
	nonce, err := sAccount.Nonce(context, rpc.BlockID{Tag: "latest"}, sAccount.AccountAddress)
	if err != nil {
		return "", err
	}

	// Here we are setting the maxFee
	maxfee, err := utils.HexToFelt("0x574fbde6000")
	if err != nil {
		return "", err
	}

	// Building the InvokeTx struct
	InvokeTx := rpc.InvokeTxnV1{
		MaxFee:        maxfee,
		Version:       rpc.TransactionV1,
		Nonce:         nonce,
		Type:          rpc.TransactionType_Invoke,
		SenderAddress: sAccount.AccountAddress,
	}

	// Building the functionCall struct, where :
	FnCall := rpc.FunctionCall{
		ContractAddress:    contractFelt,                              //contractAddress is the contract that we want to call
		EntryPointSelector: utils.GetSelectorFromNameFelt("transfer"), //this is the function that we want to call
		Calldata:           callData,
	}

	// Building the Calldata with the help of FmtCalldata where we pass in the FnCall struct along with the Cairo version
	InvokeTx.Calldata, err = sAccount.FmtCalldata([]rpc.FunctionCall{FnCall}, accountContractVersion)
	if err != nil {
		return "", err
	}

	// Signing of the transaction that is done by the account
	err = sAccount.SignInvokeTransaction(context, &InvokeTx)
	if err != nil {
		return "", err
	}

	// After the signing we finally call the AddInvokeTransaction in order to invoke the contract function
	resp, err := sAccount.AddInvokeTransaction(context, InvokeTx)
	if err != nil {
		return "", err
	}

	fmt.Println("Transaction hash response : ", resp.TransactionHash)
	_, err = sAccount.WaitForTransactionReceipt(context, resp.TransactionHash, 5)
	if err != nil {
		return "", err
	}

	// fmt.Println("Transaction hash recvResp : ")

	return resp.TransactionHash.String(), nil
}
