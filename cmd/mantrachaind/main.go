package main

import (
	"os"

	//"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/LimeChain/mantrachain/app"
)

//TODO double check and update the whole cmd functionality
func main() {
	rootCmd, _ := NewRootCmd(
	//app.Name,
	//app.AccountAddressPrefix,
	//app.DefaultNodeHome,
	//app.Name,
	//app.ModuleBasics,
	//app.New,
	// this line is used by starport scaffolding # root/arguments
	)
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
