package main

import (
	// "github.com/tolitius/stellar-mc/cmd"
	"stellar-mc/cmd"
)

// ./mc --gen-keys foo; ./mc --fund $(cat foo.pub)
func main() {

	cmd.Execute()

	// flag.StringVar(&stream, "stream", "", "stream Stellar \"ledger\", \"payments\" and \"tranasaction\" events with optional \"-s\" (seconds) and \"-c\" (cursor) subflags.\n    \texample: --stream ledger\n     \t\t--stream payments -s 42 -c now")
}
