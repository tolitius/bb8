package cmd

import (
	"log"
	"os"
	"reflect"

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
var standAloneFlag bool

// Bb8Cmd is BB-8's root command.
// Other commands are added to Bb8Cmd as subcommands.
var Bb8Cmd = &cobra.Command{
	Use:   "bb",
	Short: "cli to interact with Stellar network",
	Long:  `BB-8 is a command line interface to Stellar (https://www.stellar.org/) networks.`,
}

// AddCommands adds sub commands to Bb8Cmd
func AddCommands() {
	Bb8Cmd.AddCommand(versionCmd)
	Bb8Cmd.AddCommand(genKeysCmd)
	Bb8Cmd.AddCommand(fundCmd)
	Bb8Cmd.AddCommand(loadAccountCmd)
	Bb8Cmd.AddCommand(changeTrustCmd)
	Bb8Cmd.AddCommand(sendPaymentCmd)
	Bb8Cmd.AddCommand(streamCmd)
	Bb8Cmd.AddCommand(createAccountCmd)
	Bb8Cmd.AddCommand(manageDataCmd)
	Bb8Cmd.AddCommand(setOptionsCmd)
	Bb8Cmd.AddCommand(decodeCmd)
	Bb8Cmd.AddCommand(signTransactionCmd)
	Bb8Cmd.AddCommand(submitTransactionCmd)
}

func init() {
	withStandAlone(changeTrustCmd)
	withStandAlone(sendPaymentCmd)
	withStandAlone(setOptionsCmd)
	withStandAlone(manageDataCmd)
}

// Execute adds sub commands to Bb8Cmd and sets all the command line flags
func Execute() {

	Bb8Cmd.SilenceUsage = true
	conf = readConfig("tmp/todo")

	AddCommands()

	if c, err := Bb8Cmd.ExecuteC(); err != nil {
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

func structValues(s interface{}) []interface{} {

	v := reflect.ValueOf(s)

	values := make([]interface{}, 0)

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.IsNil() {
			values = append(values, f.Interface())
		}
	}

	return values
}

func withStandAlone(command *cobra.Command) {
	command.PersistentFlags().BoolVarP(&standAloneFlag, "sign-and-submit", "s", false,
		`sign and submit transaction. will use source account's seed to sign it

      example: send-payment -s '{"from": "seed", "to": "address", "amount": "42.0"}'
               create-account -s '{"source_account":"seed", "new_account":"address", "amount":"1.5"}'`)
}
