package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

func main() {
	configFile := flag.String("config", "config.yaml", "path to config file")
	addr := flag.String("addr", ":9200", "listen address")
	flag.Parse()

	data, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	collector, err := NewQuotaCollector(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to create collector: %v", err)
	}

	http.Handle("/metrics", promhttp.HandlerFor(collector.Registry(), promhttp.HandlerOpts{}))
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Printf("listening on %s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
