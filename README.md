# aws-quota-usage-filter

A filter plugin for [aws-quota-exporter](https://github.com/emylincon/aws_quota_exporter) that limits quota collection to only those with a CloudWatch `UsageMetric` defined.

## Problem

AWS Service Quotas returns thousands of quotas per service, but most don't have a corresponding CloudWatch usage metric. For example, EC2 returns 1,742 quotas — only 31 have usage metrics. Collecting the other 1,711 generates API calls with no observable data, inflating scrape time significantly.

## Results

| Service | Before | After | Speedup |
|---------|--------|-------|---------|
| EC2 | 94.3s | 3.1s | 30.4× |
| Lambda | 28.7s | 2.4s | 11.9× |
| RDS | 19.1s | 1.8s | 10.6× |
| **All services** | **160s** | **10.2s** | **15.7×** |

Tested across 40 AWS accounts × 4 regions.

## Installation

```bash
go install github.com/steinberg-dev/aws-quota-usage-filter@latest
```

Or build from source:

```bash
git clone https://github.com/steinberg-dev/aws-quota-usage-filter
cd aws-quota-usage-filter
go build -o aws-quota-usage-filter .
```

## Configuration

Add `collect_only_with_usage` per service in your existing config:

```yaml
quotas:
  ec2:
    collect_only_with_usage: true
  lambda:
    collect_only_with_usage: true
  rds:
    collect_only_with_usage: true
  ecs:
    collect_only_with_usage: false  # keep all quotas for this service

scrape_interval: 5m
regions:
  - us-east-1
  - eu-west-1
```

## How It Works

The filter wraps the `ListServiceQuotas` API call and discards any quota where `UsageMetric` is nil before passing results to the exporter. No changes to your existing Prometheus scrape config or dashboards required.

## License

MIT
