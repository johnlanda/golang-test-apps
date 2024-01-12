package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	addr := getEnvOrDefault("SERVER_SERVICE_URL", "localhost:8000")
	zipkinUrl := getEnvOrDefault("ZIPKIN_HTTP_URL", "http://localhost:9411/api/v2/spans")

	/* START ZIPKIN SETUP */

	// create span reporter
	reporter := httpreporter.NewReporter(zipkinUrl)
	defer reporter.Close()

	// create local service endpoint
	endpoint, err := zipkin.NewEndpoint("server", addr)
	if err != nil {
		log.Fatalf("failed to create client: %+v\n", err)
	}

	// initialize tracer
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("failed to create tracer: %+v", err)
	}

	// create server middleware
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tracer,
		zipkinhttp.TagResponseSize(true),
	)

	/* END ZIPKIN SETUP */

	r := mux.NewRouter()
	r.HandleFunc("/", serveHTTP)

	// enable the tracing middleware
	r.Use(serverMiddleware)

	log.Fatal(http.ListenAndServe(addr, r))
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print("Receiving request...")

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Hello, World.")
}

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
