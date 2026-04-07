package signal

import "testing"

func TestQualityString(t *testing.T) {
	tests := []struct {
		quality  Quality
		expected string
	}{
		{QualityExcellent, "Excellent"},
		{QualityGood, "Good"},
		{QualityFair, "Fair"},
		{QualityPoor, "Poor"},
		{QualityNone, "No Signal"},
		{Quality(99), "Unknown"},
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
		quality  Quality
		expected string
	}{
		{QualityExcellent, "★★★★★"},
		{QualityGood, "★★★★☆"},
		{QualityFair, "★★★☆☆"},
		{QualityPoor, "★★☆☆☆"},
		{QualityNone, "☆☆☆☆☆"},
		{Quality(99), "???"},
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
		metric   Metric
		expected string
	}{
		{MetricRSRP, "dBm"},
		{MetricRSRQ, "dB"},
		{MetricRSSI, "dBm"},
		{MetricSINR, "dB"},
		{Metric("UNKNOWN"), ""},
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
	rater := NewRater()

	tests := []struct {
		name     string
		rsrp     int
		expected Quality
	}{
		{"excellent signal", -80, QualityExcellent},
		{"excellent boundary", -89, QualityExcellent},
		{"good signal upper", -90, QualityGood},
		{"good signal middle", -95, QualityGood},
		{"good signal lower", -104, QualityGood},
		{"fair signal upper", -105, QualityFair},
		{"fair signal middle", -110, QualityFair},
		{"fair signal lower", -114, QualityFair},
		{"poor signal upper", -115, QualityPoor},
		{"poor signal middle", -120, QualityPoor},
		{"poor signal lower", -124, QualityPoor},
		{"no signal", -130, QualityNone},
		{"very poor signal", -140, QualityNone},
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
			if rating.Metric != MetricRSRP {
				t.Errorf("RateRSRP(%d).Metric = %v, want %v", tt.rsrp, rating.Metric, MetricRSRP)
			}
		})
	}
}

func TestRateRSRQ(t *testing.T) {
	rater := NewRater()

	tests := []struct {
		name     string
		rsrq     int
		expected Quality
	}{
		{"excellent signal", -5, QualityExcellent},
		{"excellent boundary", -9, QualityExcellent},
		{"good signal upper", -10, QualityGood},
		{"good signal middle", -12, QualityGood},
		{"good signal lower", -14, QualityGood},
		{"fair signal upper", -15, QualityFair},
		{"fair signal middle", -17, QualityFair},
		{"fair signal lower", -19, QualityFair},
		{"poor signal", -20, QualityPoor},
		{"very poor signal", -30, QualityPoor},
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
			if rating.Metric != MetricRSRQ {
				t.Errorf("RateRSRQ(%d).Metric = %v, want %v", tt.rsrq, rating.Metric, MetricRSRQ)
			}
		})
	}
}

func TestRateRSSI(t *testing.T) {
	rater := NewRater()

	tests := []struct {
		name     string
		rssi     int
		expected Quality
	}{
		{"excellent signal", -50, QualityExcellent},
		{"excellent boundary", -65, QualityExcellent},
		{"good signal upper", -70, QualityGood},
		{"good signal middle", -70, QualityGood},
		{"good signal lower", -75, QualityGood},
		{"fair signal upper", -80, QualityFair},
		{"fair signal lower", -85, QualityFair},
		{"poor signal", -90, QualityPoor},
		{"very poor signal", -100, QualityPoor},
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
			if rating.Metric != MetricRSSI {
				t.Errorf("RateRSSI(%d).Metric = %v, want %v", tt.rssi, rating.Metric, MetricRSSI)
			}
		})
	}
}

func TestRateSINR(t *testing.T) {
	rater := NewRater()

	tests := []struct {
		name     string
		sinr     int
		expected Quality
	}{
		{"excellent signal", 20, QualityExcellent},
		{"excellent boundary", 13, QualityExcellent},
		{"good signal upper", 10, QualityGood},
		{"good signal middle", 8, QualityGood},
		{"good signal lower", 6, QualityGood},
		{"fair signal upper", 5, QualityFair},
		{"fair signal middle", 3, QualityFair},
		{"fair signal lower", 0, QualityFair},
		{"poor signal", -5, QualityPoor},
		{"very poor signal", -20, QualityPoor},
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
			if rating.Metric != MetricSINR {
				t.Errorf("RateSINR(%d).Metric = %v, want %v", tt.sinr, rating.Metric, MetricSINR)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	rater := NewRater()

	tests := []struct {
		name     string
		rating   Rating
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
	customRSRP := []Threshold{
		{MinValue: -80, MaxValue: 0, Quality: QualityExcellent},
		{MinValue: -100, MaxValue: -80, Quality: QualityGood},
		{MinValue: -200, MaxValue: -100, Quality: QualityPoor},
	}

	rater := NewRaterWithThresholds(WithRSRPThresholds(customRSRP))

	rating := rater.RateRSRP(-90)
	if rating.Quality != QualityGood {
		t.Errorf("RateRSRP(-90) with custom thresholds = %v, want %v", rating.Quality, QualityGood)
	}
}

func TestWithRSRQThresholds(t *testing.T) {
	customRSRQ := []Threshold{
		{MinValue: -5, MaxValue: 20, Quality: QualityExcellent},
		{MinValue: -15, MaxValue: -5, Quality: QualityGood},
		{MinValue: -50, MaxValue: -15, Quality: QualityPoor},
	}

	rater := NewRaterWithThresholds(WithRSRQThresholds(customRSRQ))

	rating := rater.RateRSRQ(-10)
	if rating.Quality != QualityGood {
		t.Errorf("RateRSRQ(-10) with custom thresholds = %v, want %v", rating.Quality, QualityGood)
	}
}

func TestWithRSSIThresholds(t *testing.T) {
	customRSSI := []Threshold{
		{MinValue: -60, MaxValue: 0, Quality: QualityExcellent},
		{MinValue: -80, MaxValue: -60, Quality: QualityGood},
		{MinValue: -120, MaxValue: -80, Quality: QualityPoor},
	}

	rater := NewRaterWithThresholds(WithRSSIThresholds(customRSSI))

	rating := rater.RateRSSI(-70)
	if rating.Quality != QualityGood {
		t.Errorf("RateRSSI(-70) with custom thresholds = %v, want %v", rating.Quality, QualityGood)
	}
}

func TestWithSINRThresholds(t *testing.T) {
	customSINR := []Threshold{
		{MinValue: 15, MaxValue: 100, Quality: QualityExcellent},
		{MinValue: 5, MaxValue: 15, Quality: QualityGood},
		{MinValue: -100, MaxValue: 5, Quality: QualityPoor},
	}

	rater := NewRaterWithThresholds(WithSINRThresholds(customSINR))

	rating := rater.RateSINR(10)
	if rating.Quality != QualityGood {
		t.Errorf("RateSINR(10) with custom thresholds = %v, want %v", rating.Quality, QualityGood)
	}
}

func TestRateValueEdgeCases(t *testing.T) {
	rater := NewRater()

	// Test value above highest threshold (edge case)
	// RSRP: excellent is >= -89, test value above
	rsrpAbove := rater.RateRSRP(-50)
	if rsrpAbove.Quality != QualityExcellent {
		t.Errorf("RateRSRP(-50) = %v, want %v", rsrpAbove.Quality, QualityExcellent)
	}

	// Test value below lowest threshold (edge case)
	// RSRP: no signal is < -124
	rsrpBelow := rater.RateRSRP(-150)
	if rsrpBelow.Quality != QualityNone {
		t.Errorf("RateRSRP(-150) = %v, want %v", rsrpBelow.Quality, QualityNone)
	}

	// Test SINR above highest threshold
	sinrAbove := rater.RateSINR(50)
	if sinrAbove.Quality != QualityExcellent {
		t.Errorf("RateSINR(50) = %v, want %v", sinrAbove.Quality, QualityExcellent)
	}

	// Test SINR below lowest threshold
	sinrBelow := rater.RateSINR(-50)
	if sinrBelow.Quality != QualityPoor {
		t.Errorf("RateSINR(-50) = %v, want %v", sinrBelow.Quality, QualityPoor)
	}

	// Test RSRQ value above MaxValue boundary
	// RSRQ excellent MaxValue is 20, test with value above
	rsrqAbove := rater.RateRSRQ(10)
	if rsrqAbove.Quality != QualityExcellent {
		t.Errorf("RateRSRQ(10) = %v, want %v", rsrqAbove.Quality, QualityExcellent)
	}
}

func TestRateValueWithGapThresholds(t *testing.T) {
	// Create thresholds with a gap to test the fallback logic
	gapThresholds := []Threshold{
		{MinValue: -50, MaxValue: -40, Quality: QualityExcellent},
		{MinValue: -70, MaxValue: -60, Quality: QualityGood},
		{MinValue: -100, MaxValue: -90, Quality: QualityPoor},
	}

	rater := NewRaterWithThresholds(WithRSRPThresholds(gapThresholds))

	// Test value above all MaxValues
	// This triggers: if value >= thresholds[0].MaxValue
	rating := rater.RateRSRP(-35)
	if rating.Quality != QualityExcellent {
		t.Errorf(
			"RateRSRP(-35) with gap thresholds = %v, want %v",
			rating.Quality,
			QualityExcellent,
		)
	}

	// Test value in the gap between Excellent and Good (-50 to -60)
	// Falls through loop, below MaxValue, returns last threshold
	rating2 := rater.RateRSRP(-55)
	if rating2.Quality != QualityPoor {
		t.Errorf("RateRSRP(-55) with gap thresholds = %v, want %v", rating2.Quality, QualityPoor)
	}

	// Test value below all thresholds
	rating3 := rater.RateRSRP(-110)
	if rating3.Quality != QualityPoor {
		t.Errorf("RateRSRP(-110) with gap thresholds = %v, want %v", rating3.Quality, QualityPoor)
	}
}
