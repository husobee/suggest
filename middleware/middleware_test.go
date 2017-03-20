package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/husobee/suggest/response"
)

func TestMiddlewares(t *testing.T) {

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	BuildChain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("this is a panic")
	}), LoggingMiddleware, ResponseMiddleware, RecoveryMiddleware)(w, r)

	resp := w.Result()

	// we should get back a generic internal server error
	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("failed to catch the panic")
	}

	decoder := json.NewDecoder(resp.Body)
	var result = &response.Result{}
	decoder.Decode(result)

	if result.Status != http.StatusText(http.StatusInternalServerError) {
		t.Errorf("invalid status, expected OK got %s", result.Status)
	}
	if result.Message != "there was an unfortunate error" {
		t.Errorf(
			"invalid message, expected 'there was an unfortunate error' got %s",
			result.Message)
	}

}
