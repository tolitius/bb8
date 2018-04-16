package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	b "github.com/stellar/go/build"
)

var changeTrustCmd = &cobra.Command{
	Use:   "change-trust [args]",
	Short: "create, update, or delete a trustline",
	Long: `create, update, or delete a trustline.
this command takes parameters in JSON and has an optional "limit" param, setting it to "0" removes the trustline.

example: change-trust '{"source_account": "seed", "code": "XYZ", "issuer_address": "address", "limit": "42.0"}'`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		ct := &changeTrust{}
		if err := json.Unmarshal([]byte(args[0]), ct); err != nil {
			log.Fatalf("could not parse change trust transaction JSON data: %v", err)
		}

		if standAloneFlag {
			submitStandalone(conf, ct.SourceAccount, ct.makeOp())
		} else {
			if len(args) == 1 {
				encoded := makeEnvelope(conf, ct.SourceAccount, ct.makeOp())
				fmt.Print(encoded)
			} else {
				encoded := composeWithOps(args[1], ct.makeOp())
				fmt.Print(encoded)
			}
		}
	},
}

type changeTrust struct {
	SourceAccount string `json:"source_account"`
	IssuerAddress string `json:"issuer_address"`
	Issuer        string `json:"issuer"`
	Code, Limit   string
}

func (ct *changeTrust) makeOp() (muts []b.TransactionMutator) {

	// source := seedToPair(ct.SourceAccount)

	var limit = b.MaxLimit
	if ct.Limit != "" {
		limit = b.Limit(ct.Limit)
	}

	// adds backward compatibility for both
	if ct.IssuerAddress == "" {
		ct.IssuerAddress = ct.Issuer
	}

	// convert federation address to Stellar account address if needed
	ct.IssuerAddress = uniformAddress(ct.IssuerAddress)

	muts = []b.TransactionMutator{
		b.SourceAccount{AddressOrSeed: resolveAddress(ct.SourceAccount)},
		b.Trust(ct.Code, ct.IssuerAddress, limit)}

	return muts
}
