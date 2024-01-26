package pkg

import (
	"fmt"
	"testing"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/test-go/testify/require"
)

func TestDeployAccountPrecomputedAddress(t *testing.T) {

	provider, err := NewStarknetProvider("http://127.0.0.1:9944")
	require.NoError(t, err)

	accountAddress := "0x01bbdaa69ba493e429dc9efcbf560f876b445b505edd90d2577269c59ab4395d"
	_, _, privateKey := account.GetRandomKeys()

	// new account
	acnt, err := GetNewAccount(provider, accountAddress, privateKey.String())
	require.NoError(t, err)

	PubKey := GetPublickeyFromPrivateKey(privateKey.String())

	precomputedAddress, err := GetAccountPrecomputedAddress(provider, acnt, PubKey)
	require.NoError(t, err)

	fmt.Println("precomputedAddress: ", precomputedAddress)
}

func TestDeployAccount(t *testing.T) {
	accountAddress := "0x01bbdaa69ba493e429dc9efcbf560f876b445b505edd90d2577269c59ab4395d"
	// _, _, privateKey := account.GetRandomKeys()
	privateKey := "0x1e1baaf5ec023a3184b10920181c5c6dc2999e5d84500f88d93208a8202eecd"

	provider, err := NewStarknetProvider("http://127.0.0.1:9944")
	require.NoError(t, err)

	err = DeployAccount(provider, accountAddress, privateKey)
	require.NoError(t, err)
}

// 0x61bb1d625515906e2835ebc974885b4bd40cf416f1834dc05619ff7c082ffd5

func TestGetBalance(t *testing.T) {

}
