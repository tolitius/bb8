package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
)

var setOptionsCmd = &cobra.Command{
	Use:   "set-options [args]",
	Short: "set options on the account",
	Long: `set options on the account. this command takes parameters in JSON.
given a "source_account" sets options on it. supported options are:

  * inflation_destination
  * home_domain
  * master_weight
  * thresholds
  * set_flags
  * remove_flags
  * add_signer
  * remove_signer

example: set-options '{"source_account": "seed",
                       "home_domain": "stellar.org"}'

         set-options '{"source_account": "seed",
                       "home_domain": "stellar.org",
                       "max_weight": 1,
                       "inflation_destination": "GCCD6AJOYZCUAQLX32ZJF2MKFFAUJ53PVCFQI3RHWKL3V47QYE2BNAUT"}'

         set-options '{"source_account": "seed",
                       "thresholds": {"low": 1, "high": 1},
                       "set_flags": ["auth_revocable", "auth_required"]}'

         set-options '{"source_account": "seed",
                       "add_signer": {"address": "GCU2XASMVOOJCUAEPOEL7SHNIRJA3IRSDIE4UTXA4QLJHMB5BFXOLNOB",
                                      "weight": 3}}'`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		opts := &setOptions{}
		if err := json.Unmarshal([]byte(args[0]), opts); err != nil {
			log.Fatal(err)
		}

		var sourceAccount string
		if opts.SourceAccount != nil {
			sourceAccount = string(*opts.SourceAccount)
		}

		if standAloneFlag {
			submitStandalone(conf, sourceAccount, opts.makeOp())
		} else {
			if len(args) == 1 {
				encoded := makeEnvelope(conf, sourceAccount, opts.makeOp())
				fmt.Print(encoded)
			} else {
				encoded := composeWithOps(args[1], opts.makeOp())
				fmt.Print(encoded)
			}
		}
	},
}

type sourceAccount string
type setFlags []string
type clearFlags []string

var validFlags = []string{"auth_required", "auth_revocable", "auth_immutable"}

type setOptions struct {
	SourceAccount *sourceAccount   `json:"source_account"`
	HomeDomain    *b.HomeDomain    `json:"home_domain"`
	MasterWeight  *b.MasterWeight  `json:"master_weight"`
	InflationDest *b.InflationDest `json:"inflation_destination"`
	Thresholds    *b.Thresholds
	AddSigner     *b.Signer   `json:"add_signer"`
	RemoveSigner  *b.Signer   `json:"remove_signer"`
	SetFlags      *setFlags   `json:"set_flags"`
	ClearFlags    *clearFlags `json:"clear_flags"`
}

func (options *setOptions) makeSetOptionsBuilder() b.SetOptionsBuilder {

	values := structValues(*options)

	return b.SetOptions(values...)
}

func (options *setOptions) makeOp() (muts []b.TransactionMutator) {

	var optionsBuilder = options.makeSetOptionsBuilder()

	muts = []b.TransactionMutator{
		optionsBuilder}

	return muts
}

// MutateSetOptions for sourceAccount to ignore it
func (m sourceAccount) MutateSetOptions(o *xdr.SetOptionsOp) (err error) {
	return
}

// MutateSetOptions for setFlags sets the SetOptionsOp's setFlags field
func (m setFlags) MutateSetOptions(o *xdr.SetOptionsOp) (err error) {

	for _, flag := range m {

		switch flag {
		case "auth_required":
			b.SetAuthRequired().MutateSetOptions(o)
		case "auth_revocable":
			b.SetAuthRevocable().MutateSetOptions(o)
		case "auth_immutable":
			b.SetAuthImmutable().MutateSetOptions(o)
		default:
			return fmt.Errorf("unknown flag to set: \"%s\". possible flag values: %+v", flag, validFlags)
		}
	}

	return
}

// MutateSetOptions for clearFlags sets the ClearOptionsOp's clearFlags field
func (m clearFlags) MutateSetOptions(o *xdr.SetOptionsOp) (err error) {

	for _, flag := range m {

		switch flag {
		case "auth_required":
			b.ClearAuthRequired().MutateSetOptions(o)
		case "auth_revocable":
			b.ClearAuthRevocable().MutateSetOptions(o)
		case "auth_immutable":
			b.ClearAuthImmutable().MutateSetOptions(o)
		default:
			return fmt.Errorf("unknown flag to clear: \"%s\". possible flag values: %+v", flag, validFlags)
		}
	}

	return
}

// TODO: remove, once --set-options flag is out and commands are composable
func parseOptions(options string) b.SetOptionsBuilder {

	if options == "" {
		return b.SetOptionsBuilder{}
	}

	opts := &setOptions{}
	if err := json.Unmarshal([]byte(options), opts); err != nil {
		log.Fatal(err)
	}

	return opts.makeSetOptionsBuilder()
}
