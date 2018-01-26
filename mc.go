package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
)

type config struct {
	client  *horizon.Client
	network b.Network
}

func readConfig(cpath string) *config {
	//TODO: read config from cpath

	return &config{
		client:  horizon.DefaultTestNetClient,
		network: b.Network{network.TestNetworkPassphrase}}
}

func seedToPair(seed string) keypair.KP {

	kp, err := keypair.Parse(seed)
	if err != nil {
		log.Fatal(err)
	}

	return kp
}

func createNewKeys(fpath string) string {
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

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

	fmt.Fprint(fpub, pair.Address())
	fmt.Fprint(fseed, pair.Seed())

	fpub.Sync()
	fseed.Sync()

	log.Printf("keys are created and stored in: %s and %s\n", fpub.Name(), fseed.Name())

	return fpath
}

func loadAccount(stellar *horizon.Client, address, message string) horizon.Account {

	account, err := stellar.LoadAccount(address)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%saccount: %s.\nbalances: %+v\nmore details: %s",
		message, address, account.Balances, account.Links.Self.Href)

	return account
}

func fundTestAccount(stellar *horizon.Client, address string) horizon.Account {

	resp, err := http.Get(stellar.URL + "/friendbot?addr=" + address)
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

	return loadAccount(stellar, address, "successfully funded ")
}

func submitTransaction(stellar *horizon.Client, base64tx string) int32 {
	resp, err := stellar.SubmitTransaction(base64tx)
	if err != nil {
		log.Fatal(err)
	}

	return resp.Ledger
}

type newToken struct {
	symbol, issuerSeed, distributorSeed, limit string
}

func (t *newToken) issueNew(conf *config) b.Asset {

	issuer := seedToPair(t.issuerSeed)
	distributor := seedToPair(t.distributorSeed)

	asset := b.CreditAsset(t.symbol, issuer.Address())

	var limit = b.MaxLimit

	if t.limit != "" {
		limit = b.Limit(t.limit)
	}

	tx, err := b.Transaction(
		b.SourceAccount{distributor.Address()},
		b.AutoSequence{conf.client},
		conf.network,
		b.Trust(t.symbol, issuer.Address(), limit),
	)

	if err != nil {
		log.Fatal(err)
	}

	txe, err := tx.Sign(t.distributorSeed)
	if err != nil {
		log.Fatal(err)
	}

	txeB64, err := txe.Base64()

	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("tx base64: %s", txeB64)

	_, err = conf.client.SubmitTransaction(txeB64)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("submitted change of trust tx, horizon said: %s", resp)

	loadAccount(conf.client, distributor.Address(), fmt.Sprintf("issued trust for %s to ", t.symbol))

	return asset
}

// ./mc --gen-keys foo; ./mc --fund $(cat foo.pub)
func main() {
	var fund string
	var keyFpath string
	var txnToSubmit string
	var issueToken string

	flag.StringVar(&fund, "fund", "", "funds a test account. example: --fund address")
	flag.StringVar(&keyFpath, "gen-keys", "", "creates a pair of keys (in two files \"file-path\" and \"file-path.pub\"). example: --gen-keys file-path")
	flag.StringVar(&txnToSubmit, "submit-tx", "", "submits a base64 encoded transaction. example: --submit-tx txn")
	flag.StringVar(&issueToken, "issue-new-token", "", "issue new token. example: --issue-new-token token issuer-seed distributor-seed [limit]")

	flag.Parse()

	conf := readConfig("/tmp/todo")

	switch {
	case fund != "":
		fundTestAccount(conf.client, fund)
	case keyFpath != "":
		createNewKeys(keyFpath)
	case txnToSubmit != "":
		submitTransaction(conf.client, txnToSubmit)
	case issueToken != "":

		switch nargs := len(os.Args); nargs {
		case 5:
			(&newToken{
				symbol:          os.Args[2],
				issuerSeed:      os.Args[3],
				distributorSeed: os.Args[4]}).issueNew(conf)
		case 6:
			(&newToken{
				symbol:          os.Args[2],
				issuerSeed:      os.Args[3],
				distributorSeed: os.Args[4],
				limit:           os.Args[5]}).issueNew(conf)
		default:
			log.Fatalf("usage: --issue-new-token token issuer-seed distributor-seed [limit]\nthe arguments given are: %+v", os.Args[1:])
		}

	default:
		flag.PrintDefaults()
	}
}
