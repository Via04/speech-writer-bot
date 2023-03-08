package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Download file to specified path from specified url
func download(ctx context.Context, fpath string, url string) error {
	contn, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer contn.Close()
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", res.StatusCode)
	}
	_, err = io.Copy(contn, res.Body)
	if err != nil {
		return err
	}
	return nil
}

// get filename from full filename e.g. test.txt returns test
func getNameWExt(fname string, ext string) string {
	if ext[0] != '.' {
		// if ext provided without trailing '.', we should manually deal with it while getting name
		return fname[:len(fname)-len(ext)-1]
	}
	return fname[:len(fname)-len(ext)]
}
