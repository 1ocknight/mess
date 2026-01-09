package utils

import (
	"io"
	"net/http"
	"testing"
)

func SendFilePut(t *testing.T, url string, body io.Reader) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		t.Fatalf("failed to create upload request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected upload status: %d", resp.StatusCode)
	}
}
