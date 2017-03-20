package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
	"github.com/husobee/suggest/data"
	"github.com/husobee/suggest/middleware"
	"github.com/husobee/suggest/response"
)

// GetHandler - Retrieval of term
func GetHandler(w http.ResponseWriter, r *http.Request) {
	// grab the key from the request's query string
	var key = r.FormValue("key")
	// perform a data retrieval from the trie
	payload, err := data.Retrieve(key)
	if err != nil {
		// handle error (gracefully :))
		var result = response.Result{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: "failure in retrieving results",
		}
		glog.Errorf("Error retrieving results: %v", err)
		// set result to request context
		write(r, http.StatusInternalServerError, result)
		return
	}

	var result = response.GetTermResult{
		Result: response.Result{
			Status:  http.StatusText(http.StatusOK),
			Message: "successful in retrieving results",
		},
		Payload: payload,
	}
	// all is good, return with good status code and result
	write(r, http.StatusOK, result)
}

// write - helper to get the result on the context for the response middleware
func write(r *http.Request, status int, result interface{}) {
	*r = *r.WithContext(context.WithValue(r.Context(), middleware.StatusCodeKey, status))
	*r = *r.WithContext(context.WithValue(r.Context(), middleware.ResponseStructKey, result))
}

// PostHandler - Insertion of term
func PostHandler(w http.ResponseWriter, r *http.Request) {
	// defer the close of the post body, if post body is not nil
	// that is always forgotten about.
	if glog.V(1) {
		glog.Info("starting post handler")
	}
	defer func() {
		if r.Body != nil {
			if glog.V(1) {
				glog.Info("closing request body")
			}
			r.Body.Close()
		}
	}()

	// create a new decoder based on the request body
	var decoder = json.NewDecoder(r.Body)

	// create a new data.Term pointer to deserialize to
	var term = new(data.Term)
	if glog.V(1) {
		glog.Info("starting decode of post body")
	}
	// decode the term (Term is just key and value pair)
	if err := decoder.Decode(&term); err != nil {
		glog.Errorf(
			"failed to decode request body: err=%v",
			err)
	}
	if glog.V(1) {
		glog.Infof("decode of post body complete: %v", term)
	}
	// perform insertion into tree
	data.Insert(term.Key, term.Value)
	if glog.V(1) {
		glog.Info("insertion complete")
	}

	// create a new result
	result := response.PostTermResult{
		Result: response.Result{
			Status:  http.StatusText(http.StatusOK),
			Message: "successful insertion of term",
		},
	}

	// write response
	write(r, http.StatusOK, result)
}
