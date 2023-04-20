package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Download file to specified path from specified url
func Download(ctx context.Context, fpath string, url string) error {
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

// Get filename from full filename e.g. test.txt returns test
func GetNameNoExt(fname string) (string, error) {
	nameRunes := []rune(fname)
	var extDiv int
	for i := len(nameRunes) - 1; i >= 0; i-- {
		if nameRunes[i] == '.' {
			extDiv = i
			break
		}
		if i == 0 && nameRunes[i] != '.' {
			return "", fmt.Errorf("no extension in specified file %s", fname)
		}
	}
	return string(nameRunes[0:extDiv]), nil
}
