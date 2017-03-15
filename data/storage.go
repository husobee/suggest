package data

import (
	"os"
	"strconv"

	"github.com/armon/go-radix"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

var (
	// tree - main radix data structure
	tree = radix.New()
	// insert - channel for insertion coordination
	insert chan term
	// retrieve - channel for retrieve coordination
	retrieve chan retrieveMessage

	// insertBufferSize - buffer size for insert channel
	insertBufferSize int64 = 1
	// retrieveBufferSize - buffer size for retrieve channel
	retrieveBufferSize int64 = 1
)

// RunGardener - function that maintains operations against the trie, this
// function is the central orchestrator of the data structure
func RunGardener() chan struct{} {
	if glog.V(1) {
		glog.Info("initializing gardener")
	}

	// for all of our configurations, load from environment variables
	for _, v := range []struct {
		I *int64
		S string
	}{
		{
			&insertBufferSize,
			"TRIE_INSERTION_BUFFER",
		},
		{
			&retrieveBufferSize,
			"TRIE_RETRIEVE_BUFFER",
		},
	} {
		if bufferSize := os.Getenv(v.S); bufferSize != "" {
			if glog.V(1) {
				glog.Infof("populating %s from environment: %s", v.S, bufferSize)
			}
			var err error
			if *v.I, err = strconv.ParseInt(bufferSize, 10, 64); err != nil {
				glog.Fatalf("invalid %s env variable: %s, %v", v.S, bufferSize, err)
			}
			if glog.V(1) {
				glog.Infof("pulled, and parsed %s from environment: %d", v.S, *v.I)
			}
		}
	}

	// setting up buffered channels for insert and retrieve
	insert = make(chan term, insertBufferSize)
	retrieve = make(chan retrieveMessage, retrieveBufferSize)
	// create a quit channel to signal the gardener to stop
	var quit = make(chan struct{})
	// run anonymous function which will wait for signals
	// for insertions/retrieves/deletions from the trie.
	go func() {
		for {
			select {
			case r := <-retrieve:
				// retrieve case
				var result = result{
					Terms: []Term{},
					Err:   nil,
				}
				// walk the tree beginning at the prefix we are given
				// and report back all the terms which are useful
				tree.WalkPrefix(r.Key, func(s string, v interface{}) bool {
					result.Terms = append(result.Terms, Term{s, v})
					return false
				})
				r.Result <- result
			case t := <-insert:
				// insertion case
				if _, ok := tree.Insert(t.Key, t.Value); ok {
					t.Err <- nil
				}
				t.Err <- errors.New("failed to insert into structure")
			case <-quit:
				// quit our gardener
				return
			}
		}
	}()
	return quit
}

// Term - Structure used to house a request for insertion, as well
// as a structure for the api response
type Term struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// term - internal structure used to house a request for insertion
type term struct {
	Term
	Err chan error `json:"-"`
}

// result - internal structure used to house the term results
// of a retrieve, as well as any errors that are reported back
type result struct {
	Terms []Term
	Err   error
}

// retrieveMessage - internal structure used for communicating to
// gardener for retrieval, structure includes a result channel, to
// which the result of the retrieve is supposed to go
type retrieveMessage struct {
	Key    string
	Result chan result
}

// Insert - main entrypoint for inserting a key from the radix tree
func Insert(k string, v interface{}) error {
	errChan := make(chan error)
	insert <- term{Term{k, v}, errChan}
	return <-errChan
}

// Retrieve - main entrypoint for retrieving a key from the radix tree
func Retrieve(k string) ([]Term, error) {
	resultChan := make(chan result)
	retrieve <- retrieveMessage{k, resultChan}
	result := <-resultChan
	return result.Terms, result.Err
}
