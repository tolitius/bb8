package main

import (
	"flag"
	"fmt"

	// "github.com/tolitius/stellar-mc/cmd"
	"stellar-mc/cmd"

	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

type config struct {
	client  *horizon.Client
	network b.Network
}

// ./mc --gen-keys foo; ./mc --fund $(cat foo.pub)
func main() {

	cmd.Execute()

	var txOptions string
	var stream string

	flag.StringVar(&txOptions, "tx-options", "", "add one or more transaction options.\n    \texample: --tx-options '{\"home-domain\": \"stellar.org\", \"max-weight\": 1, \"inflation-destination\": \"address\"}'")
	flag.StringVar(&stream, "stream", "", "stream Stellar \"ledger\", \"payments\" and \"tranasaction\" events with optional \"-s\" (seconds) and \"-c\" (cursor) subflags.\n    \texample: --stream ledger\n     \t\t--stream payments -s 42 -c now")

	flag.Parse()

	switch {
	case stream != "":
		// opts := streamFlag.Parse(os.Args[3:])
		// fmt.Printf("stream opts: %+v", opts)
		fmt.Println("here")
	default:
		if txOptions != "" {
			fmt.Errorf("\"--tx-options\" can't be used by itself, it is an additional flag that should be used with other flags that build transactions: i.e. \"--send-payment ... --tx-options ...\" or \"--change-trust ... --tx-options ...\"")
			// } else {
			// 	flag.PrintDefaults()
			//
		}
	}
}
