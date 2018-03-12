package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	client  *horizon.Client
	network b.Network
}

var conf *config
var standAloneFlag bool
var networkName string
var horizonName string
var horizonURL string

// bb8Cmd is BB-8's root command.
// Other commands are added to bb8Cmd as subcommands.
var bb8Cmd = &cobra.Command{
	Use:   "bb",
	Short: "cli to interact with Stellar network",
	Long:  `BB-8 is a command line interface to Stellar (https://www.stellar.org/) networks.`,
}

// AddCommands adds sub commands to bb8Cmd
func AddCommands() {
	bb8Cmd.AddCommand(versionCmd)
	bb8Cmd.AddCommand(genKeysCmd)
	bb8Cmd.AddCommand(fundCmd)
	bb8Cmd.AddCommand(loadAccountCmd)
	bb8Cmd.AddCommand(changeTrustCmd)
	bb8Cmd.AddCommand(sendPaymentCmd)
	bb8Cmd.AddCommand(streamCmd)
	bb8Cmd.AddCommand(createAccountCmd)
	bb8Cmd.AddCommand(manageDataCmd)
	bb8Cmd.AddCommand(setOptionsCmd)
	bb8Cmd.AddCommand(decodeCmd)
	bb8Cmd.AddCommand(signTransactionCmd)
	bb8Cmd.AddCommand(submitTransactionCmd)
	bb8Cmd.AddCommand(accountMergeCmd)
}

func init() {

	bb8Cmd.PersistentFlags().StringVarP(&networkName, "network", "n", "", "network name from the BB8 config file. by default BB8 would look in \"home.dir/.bb8/bb8.json\" file")
	bb8Cmd.PersistentFlags().StringVar(&horizonName, "horizon", "", "horizon server name from the BB8 config file. by default BB8 would look in \"home.dir/.bb8/bb8.json\" file")
	bb8Cmd.PersistentFlags().StringVar(&horizonURL, "horizon-url", "", "horizon server URL to use for \"this\" transaction")

	withStandAlone(changeTrustCmd)
	withStandAlone(sendPaymentCmd)
	withStandAlone(setOptionsCmd)
	withStandAlone(manageDataCmd)
	withStandAlone(createAccountCmd)
	withStandAlone(accountMergeCmd)
}

// Execute adds sub commands to bb8Cmd and sets all the command line flags
func Execute() {

	bb8Cmd.SilenceUsage = true

	AddCommands()

	cobra.OnInitialize(readConfig)

	if c, err := bb8Cmd.ExecuteC(); err != nil {
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

func readConfig() {

	viper.SetConfigName("bb8")
	viper.AddConfigPath("$HOME/.bb8")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("running with defaults: no config file found.")
	}

	// read url and passphrase from config
	passphrase := viper.GetString(fmt.Sprintf("network.%s.passphrase", networkName))
	url := viper.GetString(fmt.Sprintf("network.%s.horizon.entries.%s.url", networkName, horizonName))
	defaultHorizon := viper.GetString(fmt.Sprintf("network.%s.horizon.default", networkName))

	if url == "" {
		url = viper.GetString(fmt.Sprintf("network.%s.horizon.entries.%s.url", networkName, defaultHorizon))
	}

	// overwrite url from config with --horizon-url flag value
	if horizonURL != "" {
		url = horizonURL
	}

	if url != "" && passphrase != "" {

		conf = &config{
			client: &horizon.Client{
				URL:  url,
				HTTP: http.DefaultClient,
			},
			network: b.Network{passphrase}}

		// log.Printf("running on horizon: %s\n", url)
		return
	}

	// if can't read both url and passphrase from config, check STELLAR_NETWORK
	switch snet := getEnv("STELLAR_NETWORK", "test"); snet {
	case "public":
		url = horizon.DefaultPublicNetClient.URL
		conf = &config{
			client:  horizon.DefaultPublicNetClient,
			network: b.Network{network.PublicNetworkPassphrase}}
	case "test":
		url = horizon.DefaultTestNetClient.URL
		conf = &config{
			client:  horizon.DefaultTestNetClient,
			network: b.Network{network.TestNetworkPassphrase}}
	default:
		log.Fatalf("Unknown Stellar network: \"%s\". Stellar network is either set in BB-8 config file or by the \"STELLAR_NETWORK\" environment variable. Possible values are \"public\", \"test\". An unset \"STELLAR_NETWORK\" is treated as \"test\".", snet)
	}

	// log.Printf("running on horizon: %s\n", url)
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
