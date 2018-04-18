package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/stellar/go/keypair"
)

var suffix string

var genKeysCmd = &cobra.Command{
	Use:   "gen-keys [file-name]",
	Short: "create a pair of keys (in two files \"file-name.pub\" and \"file-name\")",
	Long: `A Stellar account is represented as a pair of keys: public (a.k.a. address) and private (a.k.a. seed).
given the file name/path gen-keys generates these pair of keys in "file-name.pub" and "file-name".`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if suffix != "" {
			address, seed := createVanityKeys(suffix)
			storeKeys(address, seed, args[0])
		} else {
			address, seed := createRandomKeys()
			storeKeys(address, seed, args[0])
		}
	},
	DisableFlagsInUseLine: true,
}

func init() {
	genKeysCmd.PersistentFlags().StringVarP(&suffix, "suffix", "s", "", "generate a pair of keys where the public key ends with a particular suffix. example: --suffix DROID")
}

func storeKeys(address, seed, fpath string) string {

	fpub, err := os.Create(fpath + ".pub")
	if err != nil {
		log.Fatal(err)
	}

	fseed, err := os.Create(fpath)
	if err != nil {
		log.Fatal(err)
	}

	defer fpub.Close()
	defer fseed.Close()

	fmt.Fprint(fpub, address)
	fmt.Fprint(fseed, seed)

	fpub.Sync()
	fseed.Sync()

	log.Printf("keys are created and stored in: %s and %s\n", fpub.Name(), fseed.Name())

	return fpath
}

func createRandomKeys() (address, seed string) {
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

	return pair.Address(), pair.Seed()
}

func createVanityKeys(suffix string) (address, seed string) {

	return "foo", "bar"
}
