package signal_test

import (
	"testing"

	signal "github.com/hugoh/cellular-signal/v2"
)

func BenchmarkRatingFormat(b *testing.B) {
	rater := signal.NewRater()
	rating := rater.RateRSRP(-92)

	b.ReportAllocs()

	for range b.N {
		_ = rating.Format("%m: %v %u (%q %s)")
	}
}

func BenchmarkRaterRateRSRP(b *testing.B) {
	rater := signal.NewRater()

	b.ReportAllocs()

	for range b.N {
		_ = rater.RateRSRP(-92)
	}
}
