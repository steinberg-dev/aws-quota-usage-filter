package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/prometheus/client_golang/prometheus"
)

type QuotaCollector struct {
	cfg      Config
	registry *prometheus.Registry
	mu       sync.RWMutex
	metrics  map[string]*prometheus.GaugeVec
}

func NewQuotaCollector(ctx context.Context, cfg Config) (*QuotaCollector, error) {
	c := &QuotaCollector{
		cfg:      cfg,
		registry: prometheus.NewRegistry(),
		metrics:  make(map[string]*prometheus.GaugeVec),
	}
	go c.scrapeLoop(ctx)
	return c, nil
}

func (c *QuotaCollector) Registry() *prometheus.Registry {
	return c.registry
}

func (c *QuotaCollector) scrapeLoop(ctx context.Context) {
	ticker := time.NewTicker(c.cfg.ScrapeInterval)
	defer ticker.Stop()
	c.scrape(ctx)
	for {
		select {
		case <-ticker.C:
			c.scrape(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (c *QuotaCollector) scrape(ctx context.Context) {
	for _, region := range c.cfg.Regions {
		for svc, scfg := range c.cfg.Quotas {
			if err := c.scrapeService(ctx, region, svc, scfg); err != nil {
				log.Printf("error scraping %s/%s: %v", region, svc, err)
			}
		}
	}
}

func (c *QuotaCollector) scrapeService(ctx context.Context, region, service string, scfg ServiceConfig) error {
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return err
	}
	client := servicequotas.NewFromConfig(awsCfg)

	paginator := servicequotas.NewListServiceQuotasPaginator(client, &servicequotas.ListServiceQuotasInput{
		ServiceCode: &service,
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, q := range page.Quotas {
			if scfg.CollectOnlyWithUsage && q.UsageMetric == nil {
				continue
			}
			c.recordQuota(region, service, *q.QuotaName, *q.QuotaCode, *q.Value)
		}
	}
	return nil
}

func (c *QuotaCollector) recordQuota(region, service, name, code string, value float64) {
	key := service + "_" + code
	c.mu.Lock()
	g, ok := c.metrics[key]
	if !ok {
		g = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "aws_quota_limit",
			Help: "AWS service quota limit",
			ConstLabels: prometheus.Labels{"service": service, "quota_code": code, "quota_name": name},
		}, []string{"region"})
		c.registry.MustRegister(g)
		c.metrics[key] = g
	}
	c.mu.Unlock()
	g.WithLabelValues(region).Set(value)
}
