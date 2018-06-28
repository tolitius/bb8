package cmd

import (
	"context"
	"fmt"
	"log"
	"time"
	"strconv"

	"github.com/spf13/cobra"
	h "github.com/stellar/go/clients/horizon"
)

var ledger bool
var operation bool
var payments string
var transactions string
var seconds int
var cursor string
var startDateString string

var ledgerHandler = func(l h.Ledger) {
	fmt.Println(toJSON(l))
}

var operationHandler = func(l h.Operation) {
	fmt.Println(toJSON(l))

}

var transactionHandler = func(l h.Transaction) {
	fmt.Println(toJSON(l))
}

var paymentHandler = func(l h.Payment) {
	fmt.Println(toJSON(l))
}

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "stream \"ledger\", \"payment\", \"tranasaction\" and \"operations\" events",
	Long: `stream "ledger", "payment", "tranasaction" and "operations" events.
events are streamed in JSON and will do so forever or for a period of time specified by a --seconds flag.

example: stream --ledger
         stream -t GCYQSB3UQDSISB5LKAL2OEVLAYJNIR7LFVYDNKRMLWQKDCBX4PU3Z6JP --seconds 42 --cursor now
         stream -p luke_skywalker*scoutship.com -s 42
         stream -o --from 2018-02-12T11:45:26Z`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if seconds != 0 {
			go func() {
				time.Sleep(time.Duration(seconds) * time.Second)
				cancel()
			}()
		}

		if !ledger && !operation && payments == "" && transactions == "" {
			log.Fatal("stream command needs one of the following flags: --ledger/-l, --payments/-p, --transactions/-t --operations/-o")
		}

		var cursor = h.Cursor(cursor)

		if ledger {
			err := conf.client.StreamLedgers(ctx, &cursor, ledgerHandler)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		if operation {

			if startDateString != "" {
				op, err := getOpForCursor("now")
				if err != nil {
					log.Fatal(err)
				}

				lastPagingToken, err := strconv.ParseInt(op.PagingToken, 10, 64)
				startDate, err := time.Parse(time.RFC3339, startDateString)
				if err != nil {
					log.Fatal(err)
				}
				n, _ := closestOffset(0, lastPagingToken, startDate, 0)
				cursor = h.Cursor(strconv.FormatInt(n, 10))

			}

			err := conf.client.StreamOperations(ctx, &cursor, operationHandler)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		if transactions != "" {
			err := conf.client.StreamTransactions(ctx, uniformAddress(transactions), &cursor, transactionHandler)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		if payments != "" {
			err := conf.client.StreamPayments(ctx, uniformAddress(payments), &cursor, paymentHandler)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

	},
}

func init() {
	streamCmd.PersistentFlags().BoolVarP(&ledger, "ledger", "l", false, "stream ledger events")
	streamCmd.PersistentFlags().BoolVarP(&operation, "operations", "o", false, "stream operation-all events")
	streamCmd.PersistentFlags().StringVarP(&transactions, "transactions", "t", "", "stream account transaction events. example: --transactions account-address")
	streamCmd.PersistentFlags().StringVarP(&payments, "payments", "p", "", "stream account payment events. example: --payments account-address")
	streamCmd.PersistentFlags().IntVarP(&seconds, "seconds", "s", 0, "number of seconds to stream events for")
	streamCmd.PersistentFlags().StringVarP(&cursor, "cursor", "c", "", "a paging token, specifying where to start returning records from. When streaming this can be set to \"now\" to stream object created since your request time. examples: -c 8589934592, -c now")
	streamCmd.PersistentFlags().StringVarP(&startDateString, "from", "f", "", "a start date (RFC3339), specifying where to start returning records from when streaming operations. Override cursor setting. examples: --from 2018-02-12T11:45:26Z")
}
