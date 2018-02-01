package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"stellar-mc/cmd"

	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/network"
)

type config struct {
	client  *horizon.Client
	network b.Network
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func readConfig(cpath string) *config {

	/*TODO: add custom network support
	&config{
		client: &http.Client{
			URL:  customNetworkURL
			HTTP: http.DefaultClient,
		}

		network: b.Network{customPassphrase}}
	*/

	switch snet := getEnv("STELLAR_NETWORK", "test"); snet {
	case "public":
		return &config{
			client:  horizon.DefaultPublicNetClient,
			network: b.Network{network.PublicNetworkPassphrase}}
	case "test":
		return &config{
			client:  horizon.DefaultTestNetClient,
			network: b.Network{network.TestNetworkPassphrase}}
	default:
		log.Fatalf("Unknown Stellar network: \"%s\". Stellar network is set by the \"STELLAR_NETWORK\" environment variable. Possible values are \"public\", \"test\". An unset \"STELLAR_NETWORK\" is treated as \"test\".", snet)
	}

	return nil
}

// ./mc --gen-keys foo; ./mc --fund $(cat foo.pub)
func main() {

	cmd.Execute()

	var sendPayment string
	var txOptions string
	var buildTransaction string
	var stream string

	flag.StringVar(&sendPayment, "send-payment", "", "send payment from one account to another.\n    \texample: --send-payment '{\"from\": \"seed\", \"to\": \"address\", \"token\": \"BTC\", \"amount\": \"42.0\", \"issuer\": \"address\"}'")
	flag.StringVar(&txOptions, "tx-options", "", "add one or more transaction options.\n    \texample: --tx-options '{\"home-domain\": \"stellar.org\", \"max-weight\": 1, \"inflation-destination\": \"address\"}'")
	flag.StringVar(&buildTransaction, "new-tx", "", "build and submit a new transaction. \"operations\" and \"signers\" are optional, if there are no \"signers\", the \"source-account\" seed will be used to sign this transaction.\n    \texample: --new-tx '{\"source-account\": \"address or seed\", {\"operations\": \"trust\": {\"code\": \"XYZ\", \"issuer-address\": \"address\"}}, \"signers\": [\"seed1\", \"seed2\"]}'")
	flag.StringVar(&stream, "stream", "", "stream Stellar \"ledger\", \"payments\" and \"tranasaction\" events with optional \"-s\" (seconds) and \"-c\" (cursor) subflags.\n    \texample: --stream ledger\n     \t\t--stream payments -s 42 -c now")

	flag.Parse()

	conf := readConfig("/tmp/todo")

	var txOptionsBuilder b.SetOptionsBuilder
	if txOptions != "" {
		txOptionsBuilder = parseOptions(txOptions)
	}

	switch {
	case sendPayment != "":
		payment := &tokenPayment{}
		if err := json.Unmarshal([]byte(sendPayment), payment); err != nil {
			log.Fatal(err)
		}
		tx := payment.send(conf, txOptionsBuilder)
		submitTransaction(conf.client, tx, payment.From)
	case buildTransaction != "":
		nt := &newTransaction{}
		if err := json.Unmarshal([]byte(buildTransaction), nt); err != nil {
			log.Fatal(err)
		}
		nt.Operations.SourceAccount = &b.SourceAccount{nt.SourceAccount}
		tx := nt.Operations.buildTransaction(conf, txOptionsBuilder)
		signers := nt.Signers
		if signers == nil {
			signers = []string{nt.SourceAccount}
		}
		submitTransaction(conf.client, tx, signers...)
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
