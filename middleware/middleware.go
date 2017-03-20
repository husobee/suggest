package middleware

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/husobee/suggest/response"
)

// Middleware - a type that describes a middleware, at the core of this
// implementation a middleware is merely a function that takes a handler
// function, and returns a handler function.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// BuildChain - a function that takes a handler function, a list of middlewares
// and creates a new application stack as a single http handler
func BuildChain(f http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	// if there are no more middlewares, we just return the
	// handlerfunc, as we are done recursing.
	if len(m) == 0 {
		return f
	}
	// otherwise pop the middleware from the list,
	// and call build chain recursively as it's parameter
	return m[0](BuildChain(f, m[1:cap(m)]...))
}

// RecoveryMiddleware - takes in a http.HandlerFunc, and returns a http.HandlerFunc
func RecoveryMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if glog.V(1) {
			glog.Info("starting recovery middleware")
		}
		// defer our recovery, and save face by giving back a 500
		defer func() {
			if rec := recover(); rec != nil {
				// ouch
				*r = *r.WithContext(context.WithValue(r.Context(), StatusCodeKey, http.StatusInternalServerError))
				*r = *r.WithContext(context.WithValue(r.Context(), ResponseStructKey, response.Result{
					Status:  http.StatusText(http.StatusInternalServerError),
					Message: "there was an unfortunate error",
				}))
				glog.Errorf("Application Panic: %v", rec)
			}
			if glog.V(1) {
				glog.Info("ending recovery middleware")
			}
		}()
		// call next
		f(w, r)
	}
}

// LoggingMiddleware - takes in a http.HandlerFunc, and returns a http.HandlerFunc
func LoggingMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if glog.V(1) {
			glog.Info("starting logging middleware")
		}
		// call next
		f(w, r)
		// print log message with request details
		glog.Infof("%s %s %s %v", r.RemoteAddr, r.Method, r.URL, time.Since(start))
		if glog.V(1) {
			glog.Info("ending logging middleware")
		}
	}
}

const (
	// StatusCodeKey - Context Key for passing Status Code
	StatusCodeKey = "http_status_code"
	// ResponseStructKey - Context Key for passing Response
	ResponseStructKey = "http_response_body"
)

// contentEncoder - helper encoder interface
type contentEncoder interface {
	Encode(v interface{}) error
}

// ResponseMiddleware - takes in a http.HandlerFunc, and returns a http.HandlerFunc
func ResponseMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if glog.V(1) {
			glog.Info("starting response middleware")
		}
		// validate that we are getting a request we can service:
		accept := r.Header.Get("Accept")
		// setup encoder for response encoding
		var encoder contentEncoder

		switch accept {
		case "application/xml":
			w.Header().Set("Content-Type", "application/xml")
			encoder = xml.NewEncoder(w)
		case "application/json":
			w.Header().Set("Content-Type", "application/json")
			encoder = json.NewEncoder(w)
		case "*/*":
			w.Header().Set("Content-Type", "application/json")
			encoder = json.NewEncoder(w)
		default:
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte{})
			if r.Body != nil {
				r.Body.Close()
			}
			return
		}

		// call next
		f(w, r)
		// content negotiation

		// write the header, and response
		if status, ok := r.Context().Value(StatusCodeKey).(int); ok {
			w.WriteHeader(status)
			encoder.Encode(r.Context().Value(ResponseStructKey))
			if glog.V(1) {
				glog.Info("response: %v", r.Context().Value(ResponseStructKey))
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(response.Result{
				Status:  http.StatusText(http.StatusInternalServerError),
				Message: "There was an unfortunate error",
			})
		}
		if glog.V(1) {
			glog.Info("ending response middleware")
		}
	}
}
