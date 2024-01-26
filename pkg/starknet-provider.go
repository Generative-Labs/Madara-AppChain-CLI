package pkg

import "github.com/NethermindEth/starknet.go/rpc"

func NewStarknetProvider(url string) (*rpc.Provider, error) {

	client, err := rpc.NewClient(url)
	if err != nil {
		return nil, err
	}
	provider := rpc.NewProvider(client)

	return provider, nil
}
