// Package signal provides cellular signal quality rating utilities
// for LTE/4G/5G metrics based on industry standards.
package signal

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
)

// ErrEmptyThresholds is returned when threshold slices are empty.
var ErrEmptyThresholds = errors.New("thresholds must not be empty")

// ErrUnsortedThresholds is returned when thresholds are not sorted by
// strictly descending MinValue.
var ErrUnsortedThresholds = errors.New("thresholds must be sorted by strictly descending MinValue")

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
	Value   float64
	Metric  Metric
}

// Threshold defines the lower bound of a quality level for a signal
// metric: values greater than or equal to MinValue (and below the
// next-better threshold's MinValue) get Quality. Thresholds must be
// ordered from best quality to worst, i.e. by strictly descending
// MinValue.
type Threshold struct {
	MinValue float64
	Quality  Quality
}

// Rater provides signal rating functionality.
type Rater struct {
	thresholds map[Metric][]Threshold
}

// Option configures a Rater with custom settings.
type Option func(*Rater)

// WithThresholds sets custom thresholds for a metric. The metric does
// not have to be one of the predefined constants, so custom metrics
// can be rated too.
func WithThresholds(metric Metric, thresholds []Threshold) Option {
	return func(r *Rater) {
		r.thresholds[metric] = thresholds
	}
}

// WithRSRPThresholds sets custom RSRP thresholds.
func WithRSRPThresholds(thresholds []Threshold) Option {
	return WithThresholds(MetricRSRP, thresholds)
}

// WithRSRQThresholds sets custom RSRQ thresholds.
func WithRSRQThresholds(thresholds []Threshold) Option {
	return WithThresholds(MetricRSRQ, thresholds)
}

// WithRSSIThresholds sets custom RSSI thresholds.
func WithRSSIThresholds(thresholds []Threshold) Option {
	return WithThresholds(MetricRSSI, thresholds)
}

// WithSINRThresholds sets custom SINR thresholds.
func WithSINRThresholds(thresholds []Threshold) Option {
	return WithThresholds(MetricSINR, thresholds)
}

// NewRater creates a Rater with standard industry thresholds.
func NewRater() *Rater {
	return &Rater{
		thresholds: map[Metric][]Threshold{
			MetricRSRP: defaultRSRPThresholds(),
			MetricRSRQ: defaultRSRQThresholds(),
			MetricRSSI: defaultRSSIThresholds(),
			MetricSINR: defaultSINRThresholds(),
		},
	}
}

// NewRaterWithThresholds creates a Rater with custom thresholds.
// Returns an error if any threshold slice is empty or not sorted by
// strictly descending MinValue.
func NewRaterWithThresholds(opts ...Option) (*Rater, error) {
	rater := NewRater()
	for _, opt := range opts {
		opt(rater)
	}

	for _, metric := range slices.Sorted(maps.Keys(rater.thresholds)) {
		if err := validateThresholds(rater.thresholds[metric], string(metric)); err != nil {
			return nil, err
		}
	}

	return rater, nil
}

func validateThresholds(thresholds []Threshold, metricName string) error {
	if len(thresholds) == 0 {
		return fmt.Errorf("%s: %w", metricName, ErrEmptyThresholds)
	}

	for i := 1; i < len(thresholds); i++ {
		if thresholds[i].MinValue >= thresholds[i-1].MinValue {
			return fmt.Errorf("%s: %w", metricName, ErrUnsortedThresholds)
		}
	}

	return nil
}

// Rate rates a signal value for the given metric. Metrics without
// configured thresholds rate as QualityNone.
func (r *Rater) Rate(metric Metric, value float64) Rating {
	return Rating{
		Quality: rateValue(value, r.thresholds[metric]),
		Value:   value,
		Metric:  metric,
	}
}

// RateRSRP rates an RSRP signal value.
func (r *Rater) RateRSRP(rsrp float64) Rating {
	return r.Rate(MetricRSRP, rsrp)
}

// RateRSRQ rates an RSRQ signal value.
func (r *Rater) RateRSRQ(rsrq float64) Rating {
	return r.Rate(MetricRSRQ, rsrq)
}

// RateRSSI rates an RSSI signal value.
func (r *Rater) RateRSSI(rssi float64) Rating {
	return r.Rate(MetricRSSI, rssi)
}

// RateSINR rates a SINR signal value.
func (r *Rater) RateSINR(sinr float64) Rating {
	return r.Rate(MetricSINR, sinr)
}

// Format returns a formatted string using the default format.
// The default format is "%m: %v %u (%q %s)". See FormatWith for
// details on available format verbs.
func (r *Rater) Format(rating Rating) string {
	return r.FormatWith("%m: %v %u (%q %s)", rating)
}

// FormatWith returns a formatted string for the rating using a custom format.
// Supported verbs:
//
//	%m - metric (e.g., RSRP, RSRQ, RSSI, SINR)
//	%v - value (numeric signal value)
//	%u - unit (e.g., dBm, dB)
//	%q - quality (e.g., Excellent, Good, Fair, Poor, No Signal)
//	%s - stars (visual representation like ★★★★★)
//	%% - literal percent sign
func (*Rater) FormatWith(format string, rating Rating) string {
	var builder strings.Builder
	builder.Grow(len(format) * formatGrowthFactor)

	for idx := 0; idx < len(format); idx++ {
		if format[idx] == '%' && idx+1 < len(format) {
			appendVerb(&builder, format[idx+1], rating)
			idx++
		} else {
			builder.WriteByte(format[idx])
		}
	}

	return builder.String()
}

func appendVerb(builder *strings.Builder, verb byte, rating Rating) {
	switch verb {
	case 'm':
		builder.WriteString(rating.Metric.String())
	case 'v':
		builder.WriteString(strconv.FormatFloat(rating.Value, 'f', -1, 64))
	case 'u':
		builder.WriteString(rating.Metric.Unit())
	case 'q':
		builder.WriteString(rating.Quality.String())
	case 's':
		builder.WriteString(rating.Quality.Stars())
	case '%':
		builder.WriteByte('%')
	default:
		builder.WriteByte('%')
		builder.WriteByte(verb)
	}
}

// rateValue determines the quality rating for a given value based on
// thresholds, which must be ordered by strictly descending MinValue.
// The first threshold whose MinValue the value reaches wins; values
// below every threshold get the last (worst) threshold's Quality.
func rateValue(value float64, thresholds []Threshold) Quality {
	if len(thresholds) == 0 {
		return QualityNone
	}

	for _, t := range thresholds {
		if value >= t.MinValue {
			return t.Quality
		}
	}

	return thresholds[len(thresholds)-1].Quality
}

const formatGrowthFactor = 2
