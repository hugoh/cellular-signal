package signal_test

import (
	"fmt"

	signal "github.com/hugoh/cellular-signal"
)

func ExampleRater_RateRSRP() {
	rater := signal.NewRater()
	rating := rater.RateRSRP(-92)
	fmt.Println(rating)
	// Output: RSRP: -92 dBm (Good ★★★★☆)
}

func ExampleRating_Format() {
	rater := signal.NewRater()
	rating := rater.RateSINR(13.5)
	fmt.Println(rating.Format("%m=%v%u %s"))
	// Output: SINR=13.5dB ★★★★★
}

func ExampleRater_Rate() {
	rater := signal.NewRater()
	fmt.Println(rater.Rate(signal.MetricRSRQ, -11).Quality)
	// Output: Good
}

func ExampleNewRaterWithThresholds() {
	// Each threshold is the lower bound of a quality level, ordered
	// from best quality to worst.
	thresholds := []signal.Threshold{
		{MinValue: -80, Quality: signal.QualityExcellent},
		{MinValue: -100, Quality: signal.QualityGood},
		{MinValue: -200, Quality: signal.QualityPoor},
	}

	rater, err := signal.NewRaterWithThresholds(signal.WithRSRPThresholds(thresholds))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rater.RateRSRP(-90).Quality)
	// Output: Good
}
