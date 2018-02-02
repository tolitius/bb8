package cmd

import (
	"log"
	"os"

	b "github.com/stellar/go/build"
	"github.com/stellar/go/keypair"

	"github.com/spf13/cobra"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/network"
)

type config struct {
	client  *horizon.Client
	network b.Network
}

var conf *config

// StellarCmd is Stellar Mission Control Center's root command.
// Other commands are added to StellarCmd as subcommands.
var StellarCmd = &cobra.Command{
	Use:   "stellar-mc",
	Short: "cli to interact with Stellar network",
	Long:  `stellar is a command line interface to Stellar (https://www.stellar.org/) networks.`,
}

// AddCommands adds sub commands to StellarCmd
func AddCommands() {
	StellarCmd.AddCommand(versionCmd)
	StellarCmd.AddCommand(genKeysCmd)
	StellarCmd.AddCommand(fundCmd)
	StellarCmd.AddCommand(loadAccountCmd)
	StellarCmd.AddCommand(submitTransactionCmd)
	StellarCmd.AddCommand(changeTrustCmd)
	StellarCmd.AddCommand(sendPaymentCmd)
}

// Execute adds sub commands to StellarCmd and sets all the command line flags
func Execute() {

	StellarCmd.SilenceUsage = true
	conf = readConfig("tmp/todo")

	AddCommands()

	if c, err := StellarCmd.ExecuteC(); err != nil {
		c.Println("")
		c.Println(c.UsageString())
		os.Exit(-1)
	}
}

func seedToPair(seed string) keypair.KP {

	kp, err := keypair.Parse(seed)
	if err != nil {
		log.Fatal(err)
	}

	return kp
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
