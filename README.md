# AWS Quota Usage Filter

Filters AWS service quotas to only those with CloudWatch usage metrics. 
Reduces EC2 metrics from 1,742 to 31 — 97% scrape time reduction.

## Quick Start

```bash
go install github.com/steinberg-dev/aws-quota-usage-filter@latest
aws-quota-usage-filter --config=config.yaml --port=9200
```

Compatible with Prometheus exporters like [redis_exporter](https://github.com/oliver006/redis_exporter). 
Point your scraper at `:9200/metrics`.

## Config

See [config.example.yaml](config.example.yaml) for per-service filtering.
