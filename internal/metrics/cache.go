package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// These are attributed to be directly pulled from the memory cache that stores all the release relevant data.
var (
	TotalRepositories = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cache_repositories_total",
		Help: "Total number of repositories",
	})
	TotalReleases = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cache_releases_total",
		Help: "Total number of releases",
	})
	TotalAssets = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cache_assets_total",
		Help: "Total number of assets",
	})

	TotalDownloads = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cache_downloads_total",
		Help: "Total number of downloads",
	})
	TotalRequests = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cache_requests_total",
		Help: "Total number of requests",
	})
)
