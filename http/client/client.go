package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

func main() {
	addr := getEnvOrDefault("SERVER_SERVICE_URL", "localhost:8000")
	path := getEnvOrDefault("SERVER_SERVICE_PATH", "/")
	zipkinUrl := getEnvOrDefault("ZIPKIN_HTTP_URL", "http://localhost:9411/api/v2/spans")

	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := fmt.Sprintf("http://%s%s", addr, path)

	/* START ZIPKIN SETUP */

	// create span reporter
	reporter := httpreporter.NewReporter(zipkinUrl)
	defer reporter.Close()

	// create local service endpoint
	endpoint, err := zipkin.NewEndpoint("client", addr)
	if err != nil {
		log.Fatalf("failed to create client: %+v\n", err)
	}

	// initialize tracer
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("failed to create tracer: %+v", err)
	}

	// create global traced http client
	client, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("failed to create traced client: %+v", err)
	}
	/* END ZIPKIN SETUP */

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Printf("executing request to server...")

			var req *http.Request
			req, err := http.NewRequest(http.MethodGet, u, nil)
			if err != nil {
				log.Fatalf("failed to create request: %+v", err)
			}

			var res *http.Response
			res, err = client.DoWithAppSpan(req, "")
			if err != nil {
				log.Fatalf("unable to perform http request: %+v", err)
			}
			res.Body.Close()
		case <-interrupt:
			log.Println("interrupt")

			select {
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
