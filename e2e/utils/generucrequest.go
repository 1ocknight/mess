package utils

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

func GenericRequestWithAuth[T any, R any](t *testing.T, method string, url string, reqData *T, token string) *R {
	var result R

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	client := resty.New()

	request := client.R().
		SetContext(ctx).
		SetAuthToken(token)

	if reqData != nil {
		request.SetBody(reqData)
	}

	var resp *resty.Response
	var err error

	switch method {
	case "GET":
		resp, err = request.Get(url)
	case "POST":
		resp, err = request.Post(url)
	case "PUT":
		resp, err = request.Put(url)
	case "DELETE":
		resp, err = request.Delete(url)
	default:
		t.Fatalf("unsupported method: %s", method)
	}

	if err != nil {
		t.Fatalf("failed to request: %v", err)
	}

	if resp.StatusCode() == http.StatusNoContent {
		return nil
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		t.Fatalf("unexpected status code: %d, body: %s", resp.StatusCode(), resp.Body())
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	return &result
}
