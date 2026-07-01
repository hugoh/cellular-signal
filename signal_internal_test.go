package signal

import (
	"errors"
	"strings"
	"testing"
)

func TestRateValueEmptyThresholds(t *testing.T) {
	t.Parallel()

	result := rateValue(-50, []Threshold{})
	if result != QualityNone {
		t.Errorf("rateValue with empty thresholds = %v, want %v", result, QualityNone)
	}
}

func TestValidateThresholdsEmptyReturnsError(t *testing.T) {
	t.Parallel()

	err := validateThresholds([]Threshold{}, "RSRP")
	if err == nil {
		t.Error("validateThresholds with empty slice should return error")
	}

	if !errors.Is(err, ErrEmptyThresholds) {
		t.Errorf("validateThresholds error should wrap ErrEmptyThresholds, got: %v", err)
	}
}

func TestValidateThresholdsNilReturnsError(t *testing.T) {
	t.Parallel()

	err := validateThresholds(nil, "RSRQ")
	if err == nil {
		t.Error("validateThresholds with nil slice should return error")
	}

	if !errors.Is(err, ErrEmptyThresholds) {
		t.Errorf("validateThresholds error should wrap ErrEmptyThresholds, got: %v", err)
	}
}

func TestValidateThresholdsUnsortedReturnsError(t *testing.T) {
	t.Parallel()

	err := validateThresholds([]Threshold{
		{MinValue: -100, Quality: QualityPoor},
		{MinValue: -50, Quality: QualityExcellent},
	}, "RSRP")
	if !errors.Is(err, ErrUnsortedThresholds) {
		t.Errorf("validateThresholds error should wrap ErrUnsortedThresholds, got: %v", err)
	}
}

func TestValidateThresholdsNonEmpty(t *testing.T) {
	t.Parallel()

	err := validateThresholds([]Threshold{
		{MinValue: -100, Quality: QualityGood},
	}, "RSRP")
	if err != nil {
		t.Errorf("validateThresholds with non-empty slice should not return error: %v", err)
	}
}

// TestNewRaterWithThresholdsReportsFirstInvalidMetricSorted verifies that
// when multiple metrics are simultaneously invalid, the error names the
// metric that sorts first alphabetically, matching the deterministic
// iteration order over rater.thresholds via slices.Sorted(maps.Keys(...)).
func TestNewRaterWithThresholdsReportsFirstInvalidMetricSorted(t *testing.T) {
	t.Parallel()

	var emptyRSSI []Threshold

	var emptySINR []Threshold

	_, err := NewRaterWithThresholds(
		WithSINRThresholds(emptySINR),
		WithRSSIThresholds(emptyRSSI),
	)
	if err == nil {
		t.Fatal("NewRaterWithThresholds with multiple invalid metrics should return error")
	}

	if !strings.Contains(err.Error(), string(MetricRSSI)) {
		t.Errorf("error should report RSSI (sorts before SINR), got: %v", err)
	}
}

func FuzzRateValue(f *testing.F) {
	f.Add(-90.0)
	f.Add(0.0)
	f.Add(-200.0)

	thresholds := defaultRSRPThresholds()

	f.Fuzz(func(t *testing.T, value float64) {
		got := rateValue(value, thresholds)
		if got < QualityNone || got > QualityExcellent {
			t.Fatalf("rateValue(%v) returned out-of-range Quality %v", value, got)
		}
	})
}
