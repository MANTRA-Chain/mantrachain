package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/LimeChain/mantrachain/app"
)

func main() {
	rootCmd, _ := NewRootCmd(
	// this line is used by starport scaffolding # root/arguments
	)
	if err := svrcmd.Execute(rootCmd, "mantrachaind", app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
