package cmd

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

const (
	friendbotURL = "https://friendbot.stellar.org/?addr="
)

var fundCmd = &cobra.Command{
	Use:   "fund [address]",
	Short: "fund a test account",
	Long:  `using Stellar's Friendbot funds a Stellar test account with 10,000 lumens.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fundTestAccount(friendbotURL, args[0])
	},
	DisableFlagsInUseLine: true,
}

func fundTestAccount(url, address string) {

	resp, err := http.Get(url + uniformAddress(address))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("could not fund %s, horizon said: %s\n", address, string(body))
	}
}
