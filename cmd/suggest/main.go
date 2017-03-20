package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/husobee/suggest/data"
	"github.com/husobee/suggest/handlers"
	"github.com/husobee/suggest/middleware"
	"github.com/husobee/vestigo"
)

// config items
var (
	corsAllowOrigin      []string
	corsAllowCredentials bool
	corsExposeHeaders    []string
	corsMaxAge           time.Duration
	corsAllowHeaders     []string

	traceOn bool
	addr    string
)

func init() {
	// parse command line flags
	flag.Parse()

	// get environmental variable configurations
	corsAllowOrigin = strings.Split(os.Getenv("CORS_ALLOW_ORIGIN"), ",")
	corsExposeHeaders = strings.Split(os.Getenv("CORS_EXPOSE_HEADERS"), ",")
	corsMaxAge, _ = time.ParseDuration(os.Getenv("CORS_MAX_AGE"))
	corsAllowHeaders = strings.Split(os.Getenv("CORS_ALLOW_HEADERS"), ",")
	corsAllowCredentials, _ = strconv.ParseBool(os.Getenv("CORS_ALLOW_CREDENTIALS"))

	traceOn, _ = strconv.ParseBool(os.Getenv("TRACE"))

	addr = os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

}

func main() {
	// setup http router
	router := vestigo.NewRouter()

	// generic middlewares for all routes
	var middlewares = []middleware.Middleware{
		middleware.LoggingMiddleware,
		middleware.ResponseMiddleware,
		middleware.RecoveryMiddleware,
	}

	// define routes, get suggestions, post terms
	router.Get("/", middleware.BuildChain(
		handlers.GetHandler,
		middlewares...,
	))
	router.Post("/", middleware.BuildChain(
		handlers.PostHandler,
		middlewares...,
	))

	// Setting up router global  CORS policy
	router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowOrigin:      corsAllowOrigin,
		AllowCredentials: corsAllowCredentials,
		ExposeHeaders:    corsExposeHeaders,
		MaxAge:           corsMaxAge,
		AllowHeaders:     corsAllowHeaders,
	})

	if glog.V(1) {
		glog.Infof("cors configuration used: %v", map[string]interface{}{
			"AllowOrigin":      corsAllowOrigin,
			"AllowCredentials": corsAllowCredentials,
			"ExposeHeaders":    corsExposeHeaders,
			"MaxAge":           corsMaxAge,
			"AllowHeaders":     corsAllowHeaders,
		})
	}

	if traceOn {
		vestigo.AllowTrace = true
	}

	if glog.V(1) {
		glog.Info("starting gardener")
	}

	quitGardener := data.RunGardener()

	if glog.V(1) {
		glog.Infof("starting server on %s", addr)
	}
	glog.Warning(http.ListenAndServe(addr, router))
	glog.Warning("server stopping...")
	// signal gardener to quit
	quitGardener <- struct{}{}
	glog.Flush()
}
