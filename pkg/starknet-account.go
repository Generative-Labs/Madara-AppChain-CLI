package pkg

import (
	"errors"
	"math/big"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/sjxqqq/starknet-go/utils"
)

func GetNewAccount(provider *rpc.Provider, address string, privateKey string) (*account.Account, error) {
	AccountAddress, err := utils.HexToFelt(address)
	if err != nil {
		return nil, err
	}

	PubKey := GetPublickeyFromPrivateKey(privateKey)

	ks := account.NewMemKeystore()
	fakePrivKeyBI, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		return nil, errors.New("invalid private key")
	}

	ks.Put(PubKey, fakePrivKeyBI)

	// new account
	acnt, err := account.NewAccount(provider, AccountAddress, PubKey, ks)

	return acnt, err
}
