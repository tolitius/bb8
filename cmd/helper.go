package cmd

import (
	"encoding/json"
	"log"
	"os"
	"os/user"
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
