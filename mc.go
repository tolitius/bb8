package main

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"time"

	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

type ledgerStreamer struct {
	Seconds int64
	Cursor  horizon.Cursor
	handler horizon.LedgerHandler
}

func (streamer *ledgerStreamer) streamLedger(stellar *horizon.Client) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if streamer.Seconds != 0 {
		go func() {
			time.Sleep(time.Duration(streamer.Seconds) * time.Second)
			cancel()
		}()
	}

	err := stellar.StreamLedgers(ctx, &streamer.Cursor, streamer.handler)

	if err != nil {
		log.Fatal(err)
	}
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

type txOptions struct {
	HomeDomain    *b.HomeDomain    `json:"home-domain"`
	MasterWeight  *b.MasterWeight  `json:"master-weight"`
	InflationDest *b.InflationDest `json:"inflation-destination"`
	//TODO: add all transaction options
}

func parseOptions(options string) b.SetOptionsBuilder {
	topts := &txOptions{}
	if err := json.Unmarshal([]byte(options), topts); err != nil {
		log.Fatal(err)
	}

	values := structValues(*topts)
	return b.SetOptions(values...)
}
