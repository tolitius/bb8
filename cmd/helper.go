package cmd

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"os/user"
	"fmt"
	h "github.com/stellar/go/clients/horizon"
	"time"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func toJSON(foo interface{}) string {
	b, err := json.MarshalIndent(foo, "", "  ")
	if err != nil {
		log.Fatal("error:", err)
	}
	return string(b)
}

func homeDir() (string, error) {

	cuser, err := user.Current()
	if err != nil {
		return "", err
	}

	return cuser.HomeDir, nil
}

func Abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

type LightOperation struct {
	PagingToken     string `json:"paging_token"`
	CreatedAt       string `json:"created_at"`
	CreatedAtParsed time.Time
}

func parseLightOperation(operation interface{}) (op LightOperation, err error){
	b, err := json.Marshal(operation)
	err = json.Unmarshal(b, &op)
	if err != nil {
		return op, fmt.Errorf("Errror while Unmarshaling operation")
	}

	op.CreatedAtParsed, err = time.Parse(time.RFC3339, op.CreatedAt)
	if err != nil {
		return op, fmt.Errorf("Errror while time.Parse")
	}
	return op, nil
}

func getOpForCursor(cursor string) (op LightOperation, err error) {

	operations, _ := conf.client.LoadOperations(h.Limit(1), h.OrderDesc, h.Cursor(cursor))
	v := operations.Embedded.Records[0]

	op, err = parseLightOperation(v)
	if err != nil {
		return op, fmt.Errorf("Errror while parseLightOperation")
	}

	return op, nil

}

func closestOffset(offsetLow int64, offsetHigh int64, date time.Time, lastReadOffset int64) (candidate int64, err error) {
	if offsetLow > offsetHigh {
		return 0, fmt.Errorf("not found (offsetLow %d > offsetHigh %d)", offsetLow, offsetHigh)
	}

	candidate = (offsetLow + offsetHigh) / 2

	if candidate == lastReadOffset {
		// no move anymore, can't find value, return closest
		return candidate, nil
	}

	cur := strconv.FormatInt(candidate, 10)
	op, err := getOpForCursor(cur)
	if err != nil {
		return 0, fmt.Errorf("Errror while getOpForCursor")
	}
	//fmt.Printf("lookup low:%d high:%d avg:%d date[%s]\n", offsetLow, offsetHigh, candidate, op.CreatedAtParsed)

	if candidate == offsetLow || candidate == offsetHigh {
		return candidate, nil
	} else if op.CreatedAtParsed.After(date) {
		return closestOffset(offsetLow, candidate-1, date, candidate)
	} else if op.CreatedAtParsed.Before(date) {
		return closestOffset(candidate+1, offsetHigh, date, candidate)
	}

	return candidate, nil
}
