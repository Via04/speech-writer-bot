package config

import (
	"fmt"
	"os"
	"sync"
)

func getSecret() string {
	var syncOnce sync.Once
	var secret string
	return func() string {
		syncOnce.Do(func() {
			path, _ := os.Getwd()
			fmt.Fprint(os.Stdout, path)
			secretSlice, err := os.ReadFile("../../data/secret.txt")
			if err != nil {
				fmt.Fprintf(os.Stderr, "ioerr: no secret.txt")
			}
			secret = string(secretSlice)
		})
		return secret
	}()
}

var Secret = getSecret
