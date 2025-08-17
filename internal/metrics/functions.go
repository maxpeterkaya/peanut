package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func SetCounter(p prometheus.Gauge, v float64) {
	p.Set(v)
}

func AddCounter(p prometheus.Gauge, v float64) {
	p.Add(v)
}

func IncrementCounter(p prometheus.Gauge) {
	p.Inc()
}

func DecrementCounter(p prometheus.Gauge) {
	p.Dec()
}
