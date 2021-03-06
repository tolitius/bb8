package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/text/message"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/stellar/go/keypair"
)

var suffix string
var prefix string

var genKeysCmd = &cobra.Command{
	Use:   "gen-keys [file-name]",
	Short: "create a pair of keys (in two files \"file-name.pub\" and \"file-name\")",
	Long: `A Stellar account is represented as a pair of keys: public (a.k.a. address) and private (a.k.a. seed).
given the file name/path gen-keys generates these pair of keys in "file-name.pub" and "file-name".`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if suffix == "" && prefix == "" {
			address, seed := createRandomKeys()
			storeKeys(address, seed, args[0])
		} else {
			address, seed := createVanityKeys(prefix, suffix)
			storeKeys(address, seed, args[0])
		}
	},
	DisableFlagsInUseLine: true,
}

func init() {
	genKeysCmd.PersistentFlags().StringVarP(&suffix, "suffix", "s", "", "generate a pair of keys where the public key ends with a particular suffix. example: --suffix DROID")
	genKeysCmd.PersistentFlags().StringVarP(&prefix, "prefix", "p", "", "generate a pair of keys where the public key ends with a particular prefix. example: --prefix BEEATE")
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

func lookForKeys(
	prefix, suffix string,
	found chan *keypair.Full,
	counter *uint64) {

	for {
		kp, err := keypair.Random()
		if err != nil {
			log.Fatalf("could not generate keys due to %v", err)
		}

		atomic.AddUint64(counter, 1)

		pub := kp.Address()
		var match bool

		if suffix != "" {
			if prefix == "" {
				match = strings.HasSuffix(pub, suffix)
			} else {
				match = strings.HasPrefix(pub[1:], prefix) && strings.HasSuffix(pub, suffix)
			}
		} else {
			match = strings.HasPrefix(pub[1:], prefix)
		}

		if match {
			found <- kp
			break
		}
	}
}

type rate struct {
	count uint64
	time  int
}

func reportFoundPair(address, prefix, suffix string, tried uint64, took int) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader([]string{"address", "prefix", "suffix", "took seconds", "went through keys"})
	table.Append([]string{address, prefix, suffix, strconv.Itoa(took), strconv.FormatUint(tried, 10)})
	table.Render()
}

func createVanityKeys(prefix, suffix string) (address, seed string) {

	prefix = strings.ToUpper(prefix)
	suffix = strings.ToUpper(suffix)

	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(runtime.NumCPU())

	var counter uint64
	every := 10 * time.Second
	found := make(chan *keypair.Full)

	log.Printf("asking %d CPU cores to find keys with \"%s\" prefix and \"%s\" suffix. stand by.\n", cores, prefix, suffix)

	for i := 0; i < cores; i++ {
		go lookForKeys(prefix, suffix, found, &counter)
	}

	ticker := time.NewTicker(every)
	hashRate := &rate{}
	p := message.NewPrinter(message.MatchLanguage("en"))

	for {
		select {
		case <-ticker.C:
			interval := int(every.Seconds())
			p.Printf("went through %d keys\t| rate %d/s\t| still looking\n", counter, (counter-hashRate.count)/uint64(interval))
			hashRate.count = counter
			hashRate.time += interval
		case pair := <-found:
			reportFoundPair(pair.Address(), prefix, suffix, counter, hashRate.time)
			ticker.Stop()
			return pair.Address(), pair.Seed()
		}
	}
}
