package pm

import (
	"net/http"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	NF_COUNTER_MAP      = make(map[string]prometheus.Counter)
	NF_COUNTERVEC_MAP   = make(map[string]*prometheus.CounterVec)
	NF_GAUGE_MAP        = make(map[string]prometheus.Gauge)
	NF_GAUGEVEC_MAP     = make(map[string]*prometheus.GaugeVec)
	NF_HISTOGRAM_MAP    = make(map[string]prometheus.Histogram)
	NF_HISTOGRAMVEC_MAP = make(map[string]*prometheus.HistogramVec)
	NF_METRIC_TYPE_MAP  = make(map[string]string)
)

var (
	servicePort = "3003"
	metricURL   = "/metrics"
)

func Init(port string, url string) {
	servicePort = port
	metricURL = url
}

func Run() {
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle(metricURL, promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+servicePort, nil))
}

func Inc(metricName string, lvs ...string) {
	if typeName, ok := NF_METRIC_TYPE_MAP[metricName]; ok {
		if typeName == "Counter" {
			NF_COUNTER_MAP[metricName].Inc()
		} else if typeName == "CounterVec" {
			NF_COUNTERVEC_MAP[metricName].WithLabelValues(lvs...).Inc()
		} else if typeName == "Gauge" {
			NF_GAUGE_MAP[metricName].Inc()
		} else if typeName == "GaugeVec" {
			NF_GAUGEVEC_MAP[metricName].WithLabelValues(lvs...).Inc()
		} else {
			log.Warning(typeName + " type not support Inc operation.")
		}
	} else {
		//log.Warning("can't not found " + metricName + " in metric map")
	}
}

func Dec(metricName string, lvs ...string) {
	if typeName, ok := NF_METRIC_TYPE_MAP[metricName]; ok {
		if typeName == "Gauge" {
			NF_GAUGE_MAP[metricName].Dec()
		} else if typeName == "GaugeVec" {
			NF_GAUGEVEC_MAP[metricName].WithLabelValues(lvs...).Dec()
		} else {
			log.Warning(typeName + " type not support Dec operation.")
		}
	} else {
		//log.Warning("can't not found " + metricName + " in metric map")
	}
}

func Add(val float64, metricName string, lvs ...string) {
	if typeName, ok := NF_METRIC_TYPE_MAP[metricName]; ok {
		if typeName == "Counter" {
			NF_COUNTER_MAP[metricName].Add(val)
		} else if typeName == "CounterVec" {
			NF_COUNTERVEC_MAP[metricName].WithLabelValues(lvs...).Add(val)
		} else if typeName == "Gauge" {
			NF_GAUGE_MAP[metricName].Add(val)
		} else if typeName == "GaugeVec" {
			NF_GAUGEVEC_MAP[metricName].WithLabelValues(lvs...).Add(val)
		} else {
			log.Warning(typeName + " type not support Add operation.")
		}
	} else {
		//log.Warning("can't not found " + metricName + " in metric map")
	}
}

func Sub(val float64, metricName string, lvs ...string) {
	if typeName, ok := NF_METRIC_TYPE_MAP[metricName]; ok {
		if typeName == "Gauge" {
			NF_GAUGE_MAP[metricName].Sub(val)
		} else if typeName == "GaugeVec" {
			NF_GAUGEVEC_MAP[metricName].WithLabelValues(lvs...).Sub(val)
		} else {
			log.Warning(typeName + " type not support Sub operation.")
		}
	} else {
		//log.Warning("can't not found " + metricName + " in metric map")
	}
}

func Set(val float64, metricName string, lvs ...string) {
	if typeName, ok := NF_METRIC_TYPE_MAP[metricName]; ok {
		if typeName == "Gauge" {
			NF_GAUGE_MAP[metricName].Set(val)
		} else if typeName == "GaugeVec" {
			NF_GAUGEVEC_MAP[metricName].WithLabelValues(lvs...).Set(val)
		} else {
			log.Warning(typeName + " type not support Set operation.")
		}
	} else {
		//log.Warning("can't not found " + metricName + " in metric map")
	}
}

func Observe(val float64, metricName string, lvs ...string) {
	if typeName, ok := NF_METRIC_TYPE_MAP[metricName]; ok {
		if typeName == "Histogram" {
			NF_HISTOGRAM_MAP[metricName].Observe(val)
		} else if typeName == "HistogramVec" {
			NF_HISTOGRAMVEC_MAP[metricName].WithLabelValues(lvs...).Observe(val)
		} else {
			log.Warning(typeName + " type not support Observe operation.")
		}
	} else {
		//log.Warning("can't not found " + metricName + " in metric map")
	}
}

func RegisterMetric(metricName string, metricType string, docString string, lables []string) {
	if metricType == "Counter" {
		NF_COUNTER_MAP[metricName] = prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: metricName,
				Help: docString,
			},
		)
		prometheus.MustRegister(NF_COUNTER_MAP[metricName])

	} else if metricType == "CounterVec" {
		NF_COUNTERVEC_MAP[metricName] = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metricName,
				Help: docString,
			},
			lables,
		)
		prometheus.MustRegister(NF_COUNTERVEC_MAP[metricName])

	} else if metricType == "Gauge" {
		NF_GAUGE_MAP[metricName] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: metricName,
				Help: docString,
			},
		)
		prometheus.MustRegister(NF_GAUGE_MAP[metricName])

	} else if metricType == "GaugeVec" {
		NF_GAUGEVEC_MAP[metricName] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: metricName,
				Help: docString,
			},
			lables,
		)
		prometheus.MustRegister(NF_GAUGEVEC_MAP[metricName])

	} else {
		log.Error("not support type: " + metricType)
		return
	}

	NF_METRIC_TYPE_MAP[metricName] = metricType
}

func RegisterHistogramMetric(metricName string, metricType string, docString string, buckets []float64) {
	if metricType == "Histogram" {
		NF_HISTOGRAM_MAP[metricName] = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    metricName,
				Help:    docString,
				Buckets: buckets,
			},
		)
		prometheus.MustRegister(NF_HISTOGRAM_MAP[metricName])

	} else {
		log.Error("not support type: " + metricType)
		return
	}

	NF_METRIC_TYPE_MAP[metricName] = metricType
}

func RegisterHistogramVecMetric(metricName string, metricType string, docString string, buckets []float64, lables []string) {
	if metricType == "HistogramVec" {
		NF_HISTOGRAMVEC_MAP[metricName] = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    metricName,
				Help:    docString,
				Buckets: buckets,
			},
			lables,
		)
		prometheus.MustRegister(NF_HISTOGRAMVEC_MAP[metricName])

	} else {
		log.Error("not support type: " + metricType)
		return
	}

	NF_METRIC_TYPE_MAP[metricName] = metricType
}
