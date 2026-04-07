package signal_test

import (
	"testing"

	signal "github.com/hugoh/cellular-signal"
)

func TestQualityString(t *testing.T) {
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
			if got := tt.quality.String(); got != tt.expected {
				t.Errorf("Quality.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestQualityStars(t *testing.T) {
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
			if got := tt.quality.Stars(); got != tt.expected {
				t.Errorf("Quality.Stars() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMetricUnit(t *testing.T) {
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
			if got := tt.metric.Unit(); got != tt.expected {
				t.Errorf("Metric.Unit() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRateRSRP(t *testing.T) {
	rater := signal.NewRater()

	tests := []struct {
		name     string
		rsrp     int
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating := rater.RateRSRP(tt.rsrp)
			if rating.Quality != tt.expected {
				t.Errorf("RateRSRP(%d) = %v, want %v", tt.rsrp, rating.Quality, tt.expected)
			}

			if rating.Value != tt.rsrp {
				t.Errorf("RateRSRP(%d).Value = %v, want %v", tt.rsrp, rating.Value, tt.rsrp)
			}

			if rating.Metric != signal.MetricRSRP {
				t.Errorf(
					"RateRSRP(%d).Metric = %v, want %v",
					tt.rsrp,
					rating.Metric,
					signal.MetricRSRP,
				)
			}
		})
	}
}

func TestRateRSRQ(t *testing.T) {
	rater := signal.NewRater()

	tests := []struct {
		name     string
		rsrq     int
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
			rating := rater.RateRSRQ(tt.rsrq)
			if rating.Quality != tt.expected {
				t.Errorf("RateRSRQ(%d) = %v, want %v", tt.rsrq, rating.Quality, tt.expected)
			}

			if rating.Value != tt.rsrq {
				t.Errorf("RateRSRQ(%d).Value = %v, want %v", tt.rsrq, rating.Value, tt.rsrq)
			}

			if rating.Metric != signal.MetricRSRQ {
				t.Errorf(
					"RateRSRQ(%d).Metric = %v, want %v",
					tt.rsrq,
					rating.Metric,
					signal.MetricRSRQ,
				)
			}
		})
	}
}

func TestRateRSSI(t *testing.T) {
	rater := signal.NewRater()

	tests := []struct {
		name     string
		rssi     int
		expected signal.Quality
	}{
		{"excellent signal", -50, signal.QualityExcellent},
		{"excellent boundary", -65, signal.QualityExcellent},
		{"good signal upper", -70, signal.QualityGood},
		{"good signal middle", -70, signal.QualityGood},
		{"good signal lower", -75, signal.QualityGood},
		{"fair signal upper", -80, signal.QualityFair},
		{"fair signal lower", -85, signal.QualityFair},
		{"poor signal", -90, signal.QualityPoor},
		{"very poor signal", -100, signal.QualityPoor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rating := rater.RateRSSI(tt.rssi)
			if rating.Quality != tt.expected {
				t.Errorf("RateRSSI(%d) = %v, want %v", tt.rssi, rating.Quality, tt.expected)
			}

			if rating.Value != tt.rssi {
				t.Errorf("RateRSSI(%d).Value = %v, want %v", tt.rssi, rating.Value, tt.rssi)
			}

			if rating.Metric != signal.MetricRSSI {
				t.Errorf(
					"RateRSSI(%d).Metric = %v, want %v",
					tt.rssi,
					rating.Metric,
					signal.MetricRSSI,
				)
			}
		})
	}
}

func TestRateSINR(t *testing.T) {
	rater := signal.NewRater()

	tests := []struct {
		name     string
		sinr     int
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
			rating := rater.RateSINR(tt.sinr)
			if rating.Quality != tt.expected {
				t.Errorf("RateSINR(%d) = %v, want %v", tt.sinr, rating.Quality, tt.expected)
			}

			if rating.Value != tt.sinr {
				t.Errorf("RateSINR(%d).Value = %v, want %v", tt.sinr, rating.Value, tt.sinr)
			}

			if rating.Metric != signal.MetricSINR {
				t.Errorf(
					"RateSINR(%d).Metric = %v, want %v",
					tt.sinr,
					rating.Metric,
					signal.MetricSINR,
				)
			}
		})
	}
}

func TestFormat(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rater.Format(tt.rating); got != tt.expected {
				t.Errorf("Format() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewRaterWithThresholds(t *testing.T) {
	customRSRP := []signal.Threshold{
		{MinValue: -80, MaxValue: 0, Quality: signal.QualityExcellent},
		{MinValue: -100, MaxValue: -80, Quality: signal.QualityGood},
		{MinValue: -200, MaxValue: -100, Quality: signal.QualityPoor},
	}

	rater := signal.NewRaterWithThresholds(signal.WithRSRPThresholds(customRSRP))

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
	customRSRQ := []signal.Threshold{
		{MinValue: -5, MaxValue: 20, Quality: signal.QualityExcellent},
		{MinValue: -15, MaxValue: -5, Quality: signal.QualityGood},
		{MinValue: -50, MaxValue: -15, Quality: signal.QualityPoor},
	}

	rater := signal.NewRaterWithThresholds(signal.WithRSRQThresholds(customRSRQ))

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
	customRSSI := []signal.Threshold{
		{MinValue: -60, MaxValue: 0, Quality: signal.QualityExcellent},
		{MinValue: -80, MaxValue: -60, Quality: signal.QualityGood},
		{MinValue: -120, MaxValue: -80, Quality: signal.QualityPoor},
	}

	rater := signal.NewRaterWithThresholds(signal.WithRSSIThresholds(customRSSI))

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
	customSINR := []signal.Threshold{
		{MinValue: 15, MaxValue: 100, Quality: signal.QualityExcellent},
		{MinValue: 5, MaxValue: 15, Quality: signal.QualityGood},
		{MinValue: -100, MaxValue: 5, Quality: signal.QualityPoor},
	}

	rater := signal.NewRaterWithThresholds(signal.WithSINRThresholds(customSINR))

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

func TestRateValueWithGapThresholds(t *testing.T) {
	// Create thresholds with a gap to test the fallback logic
	gapThresholds := []signal.Threshold{
		{MinValue: -50, MaxValue: -40, Quality: signal.QualityExcellent},
		{MinValue: -70, MaxValue: -60, Quality: signal.QualityGood},
		{MinValue: -100, MaxValue: -90, Quality: signal.QualityPoor},
	}

	rater := signal.NewRaterWithThresholds(signal.WithRSRPThresholds(gapThresholds))

	// Test value above all MaxValues
	// This triggers: if value >= thresholds[0].MaxValue
	rating := rater.RateRSRP(-35)
	if rating.Quality != signal.QualityExcellent {
		t.Errorf(
			"RateRSRP(-35) with gap thresholds = %v, want %v",
			rating.Quality,
			signal.QualityExcellent,
		)
	}

	// Test value in the gap between Excellent and Good (-50 to -60)
	// Falls through loop, below MaxValue, returns last threshold
	rating2 := rater.RateRSRP(-55)
	if rating2.Quality != signal.QualityPoor {
		t.Errorf(
			"RateRSRP(-55) with gap thresholds = %v, want %v",
			rating2.Quality,
			signal.QualityPoor,
		)
	}

	// Test value below all thresholds
	rating3 := rater.RateRSRP(-110)
	if rating3.Quality != signal.QualityPoor {
		t.Errorf(
			"RateRSRP(-110) with gap thresholds = %v, want %v",
			rating3.Quality,
			signal.QualityPoor,
		)
	}
}
