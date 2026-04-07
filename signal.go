// Package signal provides cellular signal quality rating utilities
// for LTE/4G/5G metrics based on industry standards.
package signal

import "fmt"

// Quality represents a signal quality rating level.
type Quality int

const (
	// QualityNone indicates no usable signal.
	QualityNone Quality = iota
	// QualityPoor indicates poor signal quality.
	QualityPoor
	// QualityFair indicates fair signal quality.
	QualityFair
	// QualityGood indicates good signal quality.
	QualityGood
	// QualityExcellent indicates excellent signal quality.
	QualityExcellent
)

// String returns the human-readable quality name.
func (q Quality) String() string {
	switch q {
	case QualityExcellent:
		return "Excellent"
	case QualityGood:
		return "Good"
	case QualityFair:
		return "Fair"
	case QualityPoor:
		return "Poor"
	case QualityNone:
		return "No Signal"
	default:
		return "Unknown"
	}
}

// Stars returns a visual representation of signal quality using star characters.
func (q Quality) Stars() string {
	switch q {
	case QualityExcellent:
		return "★★★★★"
	case QualityGood:
		return "★★★★☆"
	case QualityFair:
		return "★★★☆☆"
	case QualityPoor:
		return "★★☆☆☆"
	case QualityNone:
		return "☆☆☆☆☆"
	default:
		return "???"
	}
}

// Metric identifies the type of signal measurement.
type Metric string

const (
	// MetricRSRP is Reference Signal Received Power (dBm).
	MetricRSRP Metric = "RSRP"
	// MetricRSRQ is Reference Signal Received Quality (dB).
	MetricRSRQ Metric = "RSRQ"
	// MetricRSSI is Received Signal Strength Indicator (dBm).
	MetricRSSI Metric = "RSSI"
	// MetricSINR is Signal to Interference-plus-Noise Ratio (dB).
	MetricSINR Metric = "SINR"
)

// String returns the metric name.
func (m Metric) String() string {
	return string(m)
}

// Unit returns the measurement unit for the metric.
func (m Metric) Unit() string {
	switch m {
	case MetricRSRP, MetricRSSI:
		return "dBm"
	case MetricRSRQ, MetricSINR:
		return "dB"
	default:
		return ""
	}
}

// Rating contains a signal quality rating with full context.
type Rating struct {
	Quality Quality
	Value   int
	Metric  Metric
}

// Threshold defines a quality boundary for a signal metric.
type Threshold struct {
	MinValue float64
	MaxValue float64
	Quality  Quality
}

// Rater provides signal rating functionality.
type Rater struct {
	rsrpThresholds []Threshold
	rsrqThresholds []Threshold
	rssiThresholds []Threshold
	sinrThresholds []Threshold
}

// Option configures a Rater with custom settings.
type Option func(*Rater)

// WithRSRPThresholds sets custom RSRP thresholds.
func WithRSRPThresholds(thresholds []Threshold) Option {
	return func(r *Rater) {
		r.rsrpThresholds = thresholds
	}
}

// WithRSRQThresholds sets custom RSRQ thresholds.
func WithRSRQThresholds(thresholds []Threshold) Option {
	return func(r *Rater) {
		r.rsrqThresholds = thresholds
	}
}

// WithRSSIThresholds sets custom RSSI thresholds.
func WithRSSIThresholds(thresholds []Threshold) Option {
	return func(r *Rater) {
		r.rssiThresholds = thresholds
	}
}

// WithSINRThresholds sets custom SINR thresholds.
func WithSINRThresholds(thresholds []Threshold) Option {
	return func(r *Rater) {
		r.sinrThresholds = thresholds
	}
}

// NewRater creates a Rater with standard industry thresholds.
func NewRater() *Rater {
	return &Rater{
		rsrpThresholds: defaultRSRPThresholds(),
		rsrqThresholds: defaultRSRQThresholds(),
		rssiThresholds: defaultRSSIThresholds(),
		sinrThresholds: defaultSINRThresholds(),
	}
}

// NewRaterWithThresholds creates a Rater with custom thresholds.
func NewRaterWithThresholds(opts ...Option) *Rater {
	r := NewRater()
	for _, opt := range opts {
		opt(r)
	}

	return r
}

// RateRSRP rates an RSRP signal value.
func (r *Rater) RateRSRP(rsrp int) Rating {
	return Rating{
		Quality: rateValue(float64(rsrp), r.rsrpThresholds),
		Value:   rsrp,
		Metric:  MetricRSRP,
	}
}

// RateRSRQ rates an RSRQ signal value.
func (r *Rater) RateRSRQ(rsrq int) Rating {
	return Rating{
		Quality: rateValue(float64(rsrq), r.rsrqThresholds),
		Value:   rsrq,
		Metric:  MetricRSRQ,
	}
}

// RateRSSI rates an RSSI signal value.
func (r *Rater) RateRSSI(rssi int) Rating {
	return Rating{
		Quality: rateValue(float64(rssi), r.rssiThresholds),
		Value:   rssi,
		Metric:  MetricRSSI,
	}
}

// RateSINR rates a SINR signal value.
func (r *Rater) RateSINR(sinr int) Rating {
	return Rating{
		Quality: rateValue(float64(sinr), r.sinrThresholds),
		Value:   sinr,
		Metric:  MetricSINR,
	}
}

// Format returns a formatted string for the rating.
func (r *Rater) Format(rating Rating) string {
	return fmt.Sprintf("%s: %d %s (%s %s)",
		rating.Metric,
		rating.Value,
		rating.Metric.Unit(),
		rating.Quality,
		rating.Quality.Stars(),
	)
}

func rateValue(value float64, thresholds []Threshold) Quality {
	for _, t := range thresholds {
		if value >= t.MinValue && value < t.MaxValue {
			return t.Quality
		}
	}

	if value >= thresholds[0].MaxValue {
		return thresholds[0].Quality
	}

	return thresholds[len(thresholds)-1].Quality
}
