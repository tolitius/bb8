package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

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

	account, err := stellar.LoadAccount(address)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("successfully funded %s.\nbalances: %+v\nmore details: %s",
		address, account.Balances, account.Links.Self.Href)

	return account
}

func submitTransaction(stellar *horizon.Client, base64tx string) int32 {
	resp, err := stellar.SubmitTransaction(base64tx)
	if err != nil {
		log.Fatal(err)
	}

	return resp.Ledger
}

func seedToPair(seed string) keypair.KP {

	kp, err := keypair.Parse(seed)
	if err != nil {
		log.Fatal(err)
	}

	return kp
}

func issueNewToken(token, limit, issuerSeed, distributorSeed string) {

	// issuer := seedToPair(issuerSeed)
	// distributor := seedToPair(distributorSeed)

	return
}

type config struct {
	client *horizon.Client
}

func readConfig(cpath string) *config {
	//TODO: read config from cpath

	return &config{client: horizon.DefaultTestNetClient}
}

// ./mc --gen-keys foo; ./mc --fund $(cat foo.pub)
func main() {
	var fund string
	var keyFpath string
	var txnToSubmit string
	var newToken string

	flag.StringVar(&fund, "fund", "", "funds a test account. example: --fund address")
	flag.StringVar(&keyFpath, "gen-keys", "", "creates a pair of keys (in two files \"file-path\" and \"file-path.pub\"). example: --gen-keys file-path")
	flag.StringVar(&txnToSubmit, "submit-tx", "", "submits a base64 encoded transaction. example: --submit-tx txn")
	flag.StringVar(&newToken, "issue-new-token", "", "issue new token. example: --issue-new-token token issuer-key-file distributor-key-file")

	flag.Parse()

	conf := readConfig("/tmp/todo")

	switch {
	case fund != "":
		fundTestAccount(conf.client, fund)
	case keyFpath != "":
		createNewKeys(keyFpath)
	case txnToSubmit != "":
		submitTransaction(conf.client, txnToSubmit)
	default:
		flag.PrintDefaults()
	}
}
