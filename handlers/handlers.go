package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
	"github.com/husobee/suggest/data"
)

// GetHandler - Retrieval of term
func GetHandler(w http.ResponseWriter, r *http.Request) {
	// grab the key from the request's query string
	var key = r.FormValue("key")
	// perform a data retrieval from the trie
	payload, err := data.Retrieve(key)
	if err != nil {
		// handle error (gracefully :))
		var result = getTermResult{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: "failure in retrieving results",
		}
		w.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(w)
		encoder.Encode(result)
		return
	}

	var result = getTermResult{
		Status:  http.StatusText(http.StatusOK),
		Message: "successful in retrieving results",
		Payload: payload,
	}
	// all is good, return with good status code and result
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(result)
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
	result := postTermResult{
		Status:  http.StatusText(http.StatusOK),
		Message: "successful insertion of term",
	}

	// write success result to caller
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(result)
}

// postTermResult - result data structure for post term endpoint
type postTermResult struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

// getTermResult - result data structure for get term endpoint
type getTermResult struct {
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}
