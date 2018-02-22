package cmd

import (
	"encoding/json"
	"log"
	"os"
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
