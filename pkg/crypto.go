package pkg

import (
	"fmt"

	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/utils"
)

func GetPublickeyFromPrivateKey(privateKey string) string {
	privInt := utils.HexToBN(privateKey)

	pubX, _, err := curve.Curve.PrivateToPoint(privInt)
	if err != nil {
		fmt.Println("can't generate public key:", err)
		panic(err)
	}

	pubkey := utils.BigToHex(pubX)

	return pubkey
}
