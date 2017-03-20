package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bouk/monkey"
	"github.com/husobee/suggest/data"
	"github.com/husobee/suggest/middleware"
	"github.com/husobee/suggest/response"
)

func TestGetHandler(t *testing.T) {
	// Patch data Retrieve to do what we want
	monkey.Patch(data.Retrieve, func(k string) ([]data.Term, error) {
		return []data.Term{
			{Key: "hello", Value: nil},
		}, nil
	})

	r := httptest.NewRequest("GET", "/?key=hello", nil)
	w := httptest.NewRecorder()

	GetHandler(w, r)

	var result = r.Context().Value(middleware.ResponseStructKey).(response.GetTermResult)

	if result.Status != http.StatusText(http.StatusOK) {
		t.Errorf("invalid status, expected OK got %s", result.Status)
	}
	if result.Message != "successful in retrieving results" {
		t.Errorf(
			"invalid message, expected 'successful in retrieving results' got %s",
			result.Message)
	}
}

func TestGetHandlerFailure(t *testing.T) {
	// Patch data Retrieve to do what we want
	monkey.Patch(data.Retrieve, func(k string) ([]data.Term, error) {
		return []data.Term{
			{Key: "hello", Value: nil},
		}, errors.New("fail")
	})

	r := httptest.NewRequest("GET", "/?key=hello", nil)
	w := httptest.NewRecorder()

	GetHandler(w, r)

	var result = r.Context().Value(middleware.ResponseStructKey).(response.Result)

	if result.Status != http.StatusText(http.StatusInternalServerError) {
		t.Errorf("invalid status, expected %s got %s",
			http.StatusText(http.StatusInternalServerError), result.Status)
	}
	if result.Message != "failure in retrieving results" {
		t.Errorf(
			"invalid message, expected 'failure in retrieving results' got %s",
			result.Message)
	}
}

func TestPostHandler(t *testing.T) {
	// Patch data Insert/Retrieve to do what we want
	monkey.Patch(data.Insert, func(k string, v interface{}) error {
		return nil
	})

	r := httptest.NewRequest("POST", "/",
		bytes.NewBufferString(`{"key": "test", "value":null}`))

	w := httptest.NewRecorder()

	PostHandler(w, r)

	var result = r.Context().Value(middleware.ResponseStructKey).(response.PostTermResult)

	if result.Status != http.StatusText(http.StatusOK) {
		t.Errorf("invalid status, expected OK got %s", result.Status)
	}
	if result.Message != "successful insertion of term" {
		t.Errorf(
			"invalid message, expected 'successful insertion of term' got %s",
			result.Message)
	}
}
