package main

import (
	"context"
	"log"
	"time"

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
