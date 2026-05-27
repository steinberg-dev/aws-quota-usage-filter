package main

import "time"

type Config struct {
	ScrapeInterval time.Duration             `yaml:"scrape_interval"`
	Regions        []string                  `yaml:"regions"`
	Quotas         map[string]ServiceConfig  `yaml:"quotas"`
}

type ServiceConfig struct {
	CollectOnlyWithUsage bool `yaml:"collect_only_with_usage"`
}
