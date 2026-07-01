package signal_test

import (
	"testing"

	signal "github.com/hugoh/cellular-signal/v2"
)

// customRSRPThresholds returns the shared custom RSRP thresholds used to
// exercise WithRSRPThresholds across multiple tests.
func customRSRPThresholds() []signal.Threshold {
	return []signal.Threshold{
		{MinValue: -80, Quality: signal.QualityExcellent},
		{MinValue: -100, Quality: signal.QualityGood},
		{MinValue: -200, Quality: signal.QualityPoor},
	}
}

// newRaterWithCustomRSRP builds a Rater using customRSRPThresholds,
// failing the test immediately on error.
func newRaterWithCustomRSRP(t *testing.T) *signal.Rater {
	t.Helper()

	rater, err := signal.NewRaterWithThresholds(signal.WithRSRPThresholds(customRSRPThresholds()))
	if err != nil {
		t.Fatalf("NewRaterWithThresholds failed: %v", err)
	}

	return rater
}

func TestQualityString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		quality  signal.Quality
		expected string
	}{
		{signal.QualityExcellent, "Excellent"},
		{signal.QualityGood, "Good"},
		{signal.QualityFair, "Fair"},
		{signal.QualityPoor, "Poor"},
		{signal.QualityNone, "No Signal"},
		{signal.Quality(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()

			if got := tt.quality.String(); got != tt.expected {
				t.Errorf("Quality.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestQualityStars(t *testing.T) {
	t.Parallel()

	tests := []struct {
		quality  signal.Quality
		expected string
	}{
		{signal.QualityExcellent, "★★★★★"},
		{signal.QualityGood, "★★★★☆"},
		{signal.QualityFair, "★★★☆☆"},
		{signal.QualityPoor, "★★☆☆☆"},
		{signal.QualityNone, "☆☆☆☆☆"},
		{signal.Quality(99), "???"},
	}

	for _, tt := range tests {
		t.Run(tt.quality.String(), func(t *testing.T) {
			t.Parallel()

			if got := tt.quality.Stars(); got != tt.expected {
				t.Errorf("Quality.Stars() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMetricUnit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		metric   signal.Metric
		expected string
	}{
		{signal.MetricRSRP, "dBm"},
		{signal.MetricRSRQ, "dB"},
		{signal.MetricRSSI, "dBm"},
		{signal.MetricSINR, "dB"},
		{signal.Metric("UNKNOWN"), ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.metric), func(t *testing.T) {
			t.Parallel()

			if got := tt.metric.Unit(); got != tt.expected {
				t.Errorf("Metric.Unit() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRateRSRP(t *testing.T) {
	t.Parallel()

	rater := signal.NewRater()

	tests := []struct {
		name     string
		rsrp     float64
		expected signal.Quality
	}{
		{"excellent signal", -80, signal.QualityExcellent},
		{"excellent boundary", -89, signal.QualityExcellent},
		{"good signal upper", -90, signal.QualityGood},
		{"good signal middle", -95, signal.QualityGood},
		{"good signal lower", -104, signal.QualityGood},
		{"fair signal upper", -105, signal.QualityFair},
		{"fair signal middle", -110, signal.QualityFair},
		{"fair signal lower", -114, signal.QualityFair},
		{"poor signal upper", -115, signal.QualityPoor},
		{"poor signal middle", -120, signal.QualityPoor},
		{"poor signal lower", -124, signal.QualityPoor},
		{"no signal", -130, signal.QualityNone},
		{"very poor signal", -140, signal.QualityNone},
		{"fractional below boundary", -89.5, signal.QualityGood},
		{"fractional above boundary", -88.5, signal.QualityExcellent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rating := rater.RateRSRP(tt.rsrp)
			if rating.Quality != tt.expected {
				t.Errorf("RateRSRP(%v) = %v, want %v", tt.rsrp, rating.Quality, tt.expected)
			}

			if rating.Value != tt.rsrp {
				t.Errorf("RateRSRP(%v).Value = %v, want %v", tt.rsrp, rating.Value, tt.rsrp)
			}

			if rating.Metric != signal.MetricRSRP {
				t.Errorf(
					"RateRSRP(%v).Metric = %v, want %v",
					tt.rsrp,
					rating.Metric,
					signal.MetricRSRP,
				)
			}
		})
	}
}

func TestRateRSRQ(t *testing.T) {
	t.Parallel()

	rater := signal.NewRater()

	tests := []struct {
		name     string
		rsrq     float64
		expected signal.Quality
	}{
		{"excellent signal", -5, signal.QualityExcellent},
		{"excellent boundary", -9, signal.QualityExcellent},
		{"good signal upper", -10, signal.QualityGood},
		{"good signal middle", -12, signal.QualityGood},
		{"good signal lower", -14, signal.QualityGood},
		{"fair signal upper", -15, signal.QualityFair},
		{"fair signal middle", -17, signal.QualityFair},
		{"fair signal lower", -19, signal.QualityFair},
		{"poor signal", -20, signal.QualityPoor},
		{"very poor signal", -30, signal.QualityPoor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rating := rater.RateRSRQ(tt.rsrq)
			if rating.Quality != tt.expected {
				t.Errorf("RateRSRQ(%v) = %v, want %v", tt.rsrq, rating.Quality, tt.expected)
			}

			if rating.Value != tt.rsrq {
				t.Errorf("RateRSRQ(%v).Value = %v, want %v", tt.rsrq, rating.Value, tt.rsrq)
			}

			if rating.Metric != signal.MetricRSRQ {
				t.Errorf(
					"RateRSRQ(%v).Metric = %v, want %v",
					tt.rsrq,
					rating.Metric,
					signal.MetricRSRQ,
				)
			}
		})
	}
}

func TestRateRSSI(t *testing.T) {
	t.Parallel()

	rater := signal.NewRater()

	tests := []struct {
		name     string
		rssi     float64
		expected signal.Quality
	}{
		{"excellent signal", -50, signal.QualityExcellent},
		{"excellent boundary", -65, signal.QualityExcellent},
		{"good signal upper", -70, signal.QualityGood},
		{"good signal middle", -72, signal.QualityGood},
		{"good signal lower", -75, signal.QualityGood},
		{"fair signal upper", -80, signal.QualityFair},
		{"fair signal lower", -85, signal.QualityFair},
		{"poor signal", -90, signal.QualityPoor},
		{"very poor signal", -100, signal.QualityPoor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rating := rater.RateRSSI(tt.rssi)
			if rating.Quality != tt.expected {
				t.Errorf("RateRSSI(%v) = %v, want %v", tt.rssi, rating.Quality, tt.expected)
			}

			if rating.Value != tt.rssi {
				t.Errorf("RateRSSI(%v).Value = %v, want %v", tt.rssi, rating.Value, tt.rssi)
			}

			if rating.Metric != signal.MetricRSSI {
				t.Errorf(
					"RateRSSI(%v).Metric = %v, want %v",
					tt.rssi,
					rating.Metric,
					signal.MetricRSSI,
				)
			}
		})
	}
}

func TestRateSINR(t *testing.T) {
	t.Parallel()

	rater := signal.NewRater()

	tests := []struct {
		name     string
		sinr     float64
		expected signal.Quality
	}{
		{"excellent signal", 20, signal.QualityExcellent},
		{"excellent boundary", 13, signal.QualityExcellent},
		{"good signal upper", 10, signal.QualityGood},
		{"good signal middle", 8, signal.QualityGood},
		{"good signal lower", 6, signal.QualityGood},
		{"fair signal upper", 5, signal.QualityFair},
		{"fair signal middle", 3, signal.QualityFair},
		{"fair signal lower", 0, signal.QualityFair},
		{"poor signal", -5, signal.QualityPoor},
		{"very poor signal", -20, signal.QualityPoor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rating := rater.RateSINR(tt.sinr)
			if rating.Quality != tt.expected {
				t.Errorf("RateSINR(%v) = %v, want %v", tt.sinr, rating.Quality, tt.expected)
			}

			if rating.Value != tt.sinr {
				t.Errorf("RateSINR(%v).Value = %v, want %v", tt.sinr, rating.Value, tt.sinr)
			}

			if rating.Metric != signal.MetricSINR {
				t.Errorf(
					"RateSINR(%v).Metric = %v, want %v",
					tt.sinr,
					rating.Metric,
					signal.MetricSINR,
				)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()

	rater := signal.NewRater()

	tests := []struct {
		name     string
		rating   signal.Rating
		expected string
	}{
		{
			name:     "RSRP Good",
			rating:   rater.RateRSRP(-92),
			expected: "RSRP: -92 dBm (Good ★★★★☆)",
		},
		{
			name:     "SINR Excellent",
			rating:   rater.RateSINR(15),
			expected: "SINR: 15 dB (Excellent ★★★★★)",
		},
		{
			name:     "RSRQ Poor",
			rating:   rater.RateRSRQ(-22),
			expected: "RSRQ: -22 dB (Poor ★★☆☆☆)",
		},
		{
			name:     "SINR fractional value",
			rating:   rater.RateSINR(13.5),
			expected: "SINR: 13.5 dB (Excellent ★★★★★)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.rating.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewRaterWithThresholds(t *testing.T) {
	t.Parallel()

	rater := newRaterWithCustomRSRP(t)

	rating := rater.RateRSRP(-90)
	if rating.Quality != signal.QualityGood {
		t.Errorf(
			"RateRSRP(-90) with custom thresholds = %v, want %v",
			rating.Quality,
			signal.QualityGood,
		)
	}
}

func TestWithRSRQThresholds(t *testing.T) {
	t.Parallel()

	customRSRQ := []signal.Threshold{
		{MinValue: -5, Quality: signal.QualityExcellent},
		{MinValue: -15, Quality: signal.QualityGood},
		{MinValue: -50, Quality: signal.QualityPoor},
	}

	rater, err := signal.NewRaterWithThresholds(signal.WithRSRQThresholds(customRSRQ))
	if err != nil {
		t.Fatalf("NewRaterWithThresholds failed: %v", err)
	}

	rating := rater.RateRSRQ(-10)
	if rating.Quality != signal.QualityGood {
		t.Errorf(
			"RateRSRQ(-10) with custom thresholds = %v, want %v",
			rating.Quality,
			signal.QualityGood,
		)
	}
}

func TestWithRSSIThresholds(t *testing.T) {
	t.Parallel()

	customRSSI := []signal.Threshold{
		{MinValue: -60, Quality: signal.QualityExcellent},
		{MinValue: -80, Quality: signal.QualityGood},
		{MinValue: -120, Quality: signal.QualityPoor},
	}

	rater, err := signal.NewRaterWithThresholds(signal.WithRSSIThresholds(customRSSI))
	if err != nil {
		t.Fatalf("NewRaterWithThresholds failed: %v", err)
	}

	rating := rater.RateRSSI(-70)
	if rating.Quality != signal.QualityGood {
		t.Errorf(
			"RateRSSI(-70) with custom thresholds = %v, want %v",
			rating.Quality,
			signal.QualityGood,
		)
	}
}

func TestWithSINRThresholds(t *testing.T) {
	t.Parallel()

	customSINR := []signal.Threshold{
		{MinValue: 15, Quality: signal.QualityExcellent},
		{MinValue: 5, Quality: signal.QualityGood},
		{MinValue: -100, Quality: signal.QualityPoor},
	}

	rater, err := signal.NewRaterWithThresholds(signal.WithSINRThresholds(customSINR))
	if err != nil {
		t.Fatalf("NewRaterWithThresholds failed: %v", err)
	}

	rating := rater.RateSINR(10)
	if rating.Quality != signal.QualityGood {
		t.Errorf(
			"RateSINR(10) with custom thresholds = %v, want %v",
			rating.Quality,
			signal.QualityGood,
		)
	}
}

func TestRateValueEdgeCases(t *testing.T) {
	t.Parallel()

	rater := signal.NewRater()

	// Test value above highest threshold (edge case)
	// RSRP: excellent is >= -89, test value above
	rsrpAbove := rater.RateRSRP(-50)
	if rsrpAbove.Quality != signal.QualityExcellent {
		t.Errorf("RateRSRP(-50) = %v, want %v", rsrpAbove.Quality, signal.QualityExcellent)
	}

	// Test value below lowest threshold (edge case)
	// RSRP: no signal is < -124
	rsrpBelow := rater.RateRSRP(-150)
	if rsrpBelow.Quality != signal.QualityNone {
		t.Errorf("RateRSRP(-150) = %v, want %v", rsrpBelow.Quality, signal.QualityNone)
	}

	// Test SINR above highest threshold
	sinrAbove := rater.RateSINR(50)
	if sinrAbove.Quality != signal.QualityExcellent {
		t.Errorf("RateSINR(50) = %v, want %v", sinrAbove.Quality, signal.QualityExcellent)
	}

	// Test SINR below lowest threshold
	sinrBelow := rater.RateSINR(-50)
	if sinrBelow.Quality != signal.QualityPoor {
		t.Errorf("RateSINR(-50) = %v, want %v", sinrBelow.Quality, signal.QualityPoor)
	}

	// Test RSRQ value above MaxValue boundary
	// RSRQ excellent MaxValue is 20, test with value above
	rsrqAbove := rater.RateRSRQ(10)
	if rsrqAbove.Quality != signal.QualityExcellent {
		t.Errorf("RateRSRQ(10) = %v, want %v", rsrqAbove.Quality, signal.QualityExcellent)
	}
}

func TestRateValueSparseThresholds(t *testing.T) {
	t.Parallel()

	sparseThresholds := []signal.Threshold{
		{MinValue: -50, Quality: signal.QualityExcellent},
		{MinValue: -70, Quality: signal.QualityGood},
		{MinValue: -100, Quality: signal.QualityPoor},
	}

	rater, err := signal.NewRaterWithThresholds(signal.WithRSRPThresholds(sparseThresholds))
	if err != nil {
		t.Fatalf("NewRaterWithThresholds failed: %v", err)
	}

	tests := []struct {
		name     string
		rsrp     float64
		expected signal.Quality
	}{
		{"above all thresholds", -35, signal.QualityExcellent},
		{"between excellent and good", -55, signal.QualityGood},
		{"below all thresholds", -110, signal.QualityPoor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rating := rater.RateRSRP(tt.rsrp)
			if rating.Quality != tt.expected {
				t.Errorf("RateRSRP(%v) = %v, want %v", tt.rsrp, rating.Quality, tt.expected)
			}
		})
	}
}

func TestNewRaterWithThresholdsUnsortedReturnsError(t *testing.T) {
	t.Parallel()

	unsorted := []signal.Threshold{
		{MinValue: -100, Quality: signal.QualityPoor},
		{MinValue: -50, Quality: signal.QualityExcellent},
	}

	_, err := signal.NewRaterWithThresholds(signal.WithRSRPThresholds(unsorted))
	if err == nil {
		t.Error("NewRaterWithThresholds with unsorted thresholds should return error")
	}
}

func TestNewRaterWithThresholdsDuplicateMinReturnsError(t *testing.T) {
	t.Parallel()

	duplicate := []signal.Threshold{
		{MinValue: -50, Quality: signal.QualityExcellent},
		{MinValue: -50, Quality: signal.QualityGood},
	}

	_, err := signal.NewRaterWithThresholds(signal.WithSINRThresholds(duplicate))
	if err == nil {
		t.Error("NewRaterWithThresholds with duplicate MinValue should return error")
	}
}

func TestFormatWithUnknownVerb(t *testing.T) {
	t.Parallel()

	rater := signal.NewRater()
	rating := rater.RateRSRP(-92)

	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{
			name:     "unknown verb 'x'",
			format:   "Signal: %x",
			expected: "Signal: %x",
		},
		{
			name:     "unknown verb 'd'",
			format:   "%d dBm",
			expected: "%d dBm",
		},
		{
			name:     "mixed known and unknown verbs",
			format:   "%m: %x %v",
			expected: "RSRP: %x -92",
		},
		{
			name:     "literal percent",
			format:   "100%% signal",
			expected: "100% signal",
		},
		{
			name:     "literal percent before verb",
			format:   "%%m is %m",
			expected: "%m is RSRP",
		},
		{
			name:     "trailing percent",
			format:   "signal %",
			expected: "signal %",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := rating.Format(tt.format); got != tt.expected {
				t.Errorf("Format(%q) = %q, want %q", tt.format, got, tt.expected)
			}
		})
	}
}

func TestNewRaterWithThresholdsEmptyReturnsError(t *testing.T) {
	t.Parallel()

	var emptyRSRP []signal.Threshold

	_, err := signal.NewRaterWithThresholds(signal.WithRSRPThresholds(emptyRSRP))
	if err == nil {
		t.Error("NewRaterWithThresholds with empty thresholds should return error")
	}
}

func TestNewRaterWithThresholdsNilReturnsError(t *testing.T) {
	t.Parallel()

	_, err := signal.NewRaterWithThresholds(signal.WithRSRQThresholds(nil))
	if err == nil {
		t.Error("NewRaterWithThresholds with nil thresholds should return error")
	}
}

func TestNewRaterWithThresholdsEmptyRSSI(t *testing.T) {
	t.Parallel()

	var emptyRSSI []signal.Threshold

	_, err := signal.NewRaterWithThresholds(signal.WithRSSIThresholds(emptyRSSI))
	if err == nil {
		t.Error("NewRaterWithThresholds with empty RSSI thresholds should return error")
	}
}

func TestNewRaterWithThresholdsNilRSSI(t *testing.T) {
	t.Parallel()

	_, err := signal.NewRaterWithThresholds(signal.WithRSSIThresholds(nil))
	if err == nil {
		t.Error("NewRaterWithThresholds with nil RSSI thresholds should return error")
	}
}

func TestNewRaterWithThresholdsEmptySINR(t *testing.T) {
	t.Parallel()

	var emptySINR []signal.Threshold

	_, err := signal.NewRaterWithThresholds(signal.WithSINRThresholds(emptySINR))
	if err == nil {
		t.Error("NewRaterWithThresholds with empty SINR thresholds should return error")
	}
}

func TestNewRaterWithThresholdsNilSINR(t *testing.T) {
	t.Parallel()

	_, err := signal.NewRaterWithThresholds(signal.WithSINRThresholds(nil))
	if err == nil {
		t.Error("NewRaterWithThresholds with nil SINR thresholds should return error")
	}
}

func TestWithThresholdsIsolatesCallerSlice(t *testing.T) {
	t.Parallel()

	thresholds := customRSRPThresholds()

	rater, err := signal.NewRaterWithThresholds(signal.WithRSRPThresholds(thresholds))
	if err != nil {
		t.Fatalf("NewRaterWithThresholds failed: %v", err)
	}

	// Mutate the original slice after rater creation.
	thresholds[0] = signal.Threshold{MinValue: 0, Quality: signal.QualityPoor}

	rating := rater.RateRSRP(-70)
	if rating.Quality != signal.QualityExcellent {
		t.Errorf(
			"RateRSRP(-70) after mutating source slice = %v, want %v (caller mutation must not affect Rater)",
			rating.Quality,
			signal.QualityExcellent,
		)
	}
}

func TestNewRaterWithThresholdsSingleThreshold(t *testing.T) {
	t.Parallel()

	singleRSRP := []signal.Threshold{
		{MinValue: -100, Quality: signal.QualityGood},
	}

	rater, err := signal.NewRaterWithThresholds(signal.WithRSRPThresholds(singleRSRP))
	if err != nil {
		t.Fatalf("NewRaterWithThresholds failed: %v", err)
	}

	rating := rater.RateRSRP(-75)
	if rating.Quality != signal.QualityGood {
		t.Errorf(
			"RateRSRP(-75) with single threshold = %v, want %v",
			rating.Quality,
			signal.QualityGood,
		)
	}

	rating2 := rater.RateRSRP(-40)
	if rating2.Quality != signal.QualityGood {
		t.Errorf(
			"RateRSRP(-40) above single threshold = %v, want %v",
			rating2.Quality,
			signal.QualityGood,
		)
	}

	rating3 := rater.RateRSRP(-150)
	if rating3.Quality != signal.QualityGood {
		t.Errorf(
			"RateRSRP(-150) below single threshold = %v, want %v",
			rating3.Quality,
			signal.QualityGood,
		)
	}
}

func TestRateUnregisteredMetricReturnsQualityNone(t *testing.T) {
	t.Parallel()

	rater := signal.NewRater()

	rating := rater.Rate(signal.Metric("UNKNOWN"), -80)
	if rating.Quality != signal.QualityNone {
		t.Errorf("Rate with unregistered metric = %v, want %v", rating.Quality, signal.QualityNone)
	}
}

func TestNewRaterWithThresholdsLeavesOtherMetricsAtDefault(t *testing.T) {
	t.Parallel()

	rater := newRaterWithCustomRSRP(t)
	defaultRater := signal.NewRater()

	for _, tt := range []struct {
		name   string
		metric signal.Metric
		value  float64
	}{
		{"RSRQ", signal.MetricRSRQ, -11},
		{"RSSI", signal.MetricRSSI, -72},
		{"SINR", signal.MetricSINR, 8},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := rater.Rate(tt.metric, tt.value).Quality
			want := defaultRater.Rate(tt.metric, tt.value).Quality

			if got != want {
				t.Errorf(
					"Rate(%v, %v) after customizing RSRP = %v, want unchanged default %v",
					tt.metric,
					tt.value,
					got,
					want,
				)
			}
		})
	}
}

func TestWithThresholdsCustomMetric(t *testing.T) {
	t.Parallel()

	customMetric := signal.Metric("BATTERY")
	thresholds := []signal.Threshold{
		{MinValue: 80, Quality: signal.QualityExcellent},
		{MinValue: 40, Quality: signal.QualityGood},
		{MinValue: 0, Quality: signal.QualityPoor},
	}

	rater, err := signal.NewRaterWithThresholds(signal.WithThresholds(customMetric, thresholds))
	if err != nil {
		t.Fatalf("NewRaterWithThresholds failed: %v", err)
	}

	rating := rater.Rate(customMetric, 50)
	if rating.Quality != signal.QualityGood {
		t.Errorf("Rate(%v, 50) = %v, want %v", customMetric, rating.Quality, signal.QualityGood)
	}
}
