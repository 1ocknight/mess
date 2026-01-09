package utils

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func Checkfile(t *testing.T, url string, content []byte) {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("failed to GET avatar: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read avatar body: %v", err)
	}

	if !bytes.Equal(body, content) {
		t.Fatal("avatar content does not match uploaded content")
	}
}
