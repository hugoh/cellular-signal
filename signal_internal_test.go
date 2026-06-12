package signal

import (
	"errors"
	"testing"
)

func TestRateValueEmptyThresholds(t *testing.T) {
	result := rateValue(-50, []Threshold{})
	if result != QualityNone {
		t.Errorf("rateValue with empty thresholds = %v, want %v", result, QualityNone)
	}
}

func TestValidateThresholdsEmptyReturnsError(t *testing.T) {
	err := validateThresholds([]Threshold{}, "RSRP")
	if err == nil {
		t.Error("validateThresholds with empty slice should return error")
	}

	if !errors.Is(err, ErrEmptyThresholds) {
		t.Errorf("validateThresholds error should wrap ErrEmptyThresholds, got: %v", err)
	}
}

func TestValidateThresholdsNilReturnsError(t *testing.T) {
	err := validateThresholds(nil, "RSRQ")
	if err == nil {
		t.Error("validateThresholds with nil slice should return error")
	}

	if !errors.Is(err, ErrEmptyThresholds) {
		t.Errorf("validateThresholds error should wrap ErrEmptyThresholds, got: %v", err)
	}
}

func TestValidateThresholdsUnsortedReturnsError(t *testing.T) {
	err := validateThresholds([]Threshold{
		{MinValue: -100, Quality: QualityPoor},
		{MinValue: -50, Quality: QualityExcellent},
	}, "RSRP")
	if !errors.Is(err, ErrUnsortedThresholds) {
		t.Errorf("validateThresholds error should wrap ErrUnsortedThresholds, got: %v", err)
	}
}

func TestValidateThresholdsNonEmpty(t *testing.T) {
	err := validateThresholds([]Threshold{
		{MinValue: -100, Quality: QualityGood},
	}, "RSRP")
	if err != nil {
		t.Errorf("validateThresholds with non-empty slice should not return error: %v", err)
	}
}
