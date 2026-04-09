# cellular-signal

[![Go Reference](https://pkg.go.dev/badge/github.com/hugoh/cellular-signal.svg)](https://pkg.go.dev/github.com/hugoh/cellular-signal)
[![CI](https://github.com/hugoh/cellular-signal/actions/workflows/ci.yml/badge.svg)](https://github.com/hugoh/cellular-signal/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/hugoh/cellular-signal/graph/badge.svg?token=UMZMODZ5PV)](https://codecov.io/github/hugoh/cellular-signal)
[![Go Report Card](https://goreportcard.com/badge/github.com/hugoh/cellular-signal)](https://goreportcard.com/report/github.com/hugoh/cellular-signal)

Go library for rating cellular signal quality (LTE/4G/5G) based on industry standards.

## Features

- Rate RSRP, RSRQ, RSSI, and SINR signal metrics
- Based on industry-standard thresholds from telecom vendors
- Zero dependencies (pure Go stdlib)
- Fully tested
- Customizable thresholds

## Installation

```bash
go get github.com/hugoh/cellular-signal
```

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/hugoh/cellular-signal"
)

func main() {
    rater := signal.NewRater()

    // Rate individual metrics
    rsrpRating := rater.RateRSRP(-92)
    fmt.Println(rater.Format(rsrpRating))
    // Output: RSRP: -92 dBm (Good ★★★★☆)

    // Access rating details
    fmt.Printf("Quality: %s\n", rsrpRating.Quality.String())
    fmt.Printf("Metric: %s\n", rsrpRating.Metric)
    fmt.Printf("Value: %d %s\n", rsrpRating.Value, rsrpRating.Metric.Unit())
}
```

## API Reference

### Creating a Rater

```go
// Default rater with industry-standard thresholds
rater := signal.NewRater()

// Custom thresholds
customThresholds := []signal.Threshold{
    {MinValue: -80, MaxValue: 0, Quality: signal.QualityExcellent},
    {MinValue: -100, MaxValue: -80, Quality: signal.QualityGood},
    {MinValue: -200, MaxValue: -100, Quality: signal.QualityPoor},
}
rater = signal.NewRaterWithThresholds(
    signal.WithRSRPThresholds(customThresholds),
)
```

### Rating Signals

```go
// Rate RSRP (Reference Signal Received Power)
rsrpRating := rater.RateRSRP(-92)

// Rate RSRQ (Reference Signal Received Quality)
rsrqRating := rater.RateRSRQ(-11)

// Rate RSSI (Received Signal Strength Indicator)
rssiRating := rater.RateRSSI(-68)

// Rate SINR (Signal to Interference-plus-Noise Ratio)
sinrRating := rater.RateSINR(8)
```

### Quality Levels

```go
// Quality constants
signal.QualityExcellent  // "Excellent"
signal.QualityGood       // "Good"
signal.QualityFair       // "Fair"
signal.QualityPoor       // "Poor"
signal.QualityNone       // "No Signal"

// String representation
quality.String()  // Human-readable name

// Visual representation
quality.Stars()   // Star representation (★★★★★, ★★★★☆, etc.)
```

### Formatting Output

```go
rating := rater.RateRSRP(-92)
formatted := rater.Format(rating)
// Output: "RSRP: -92 dBm (Good ★★★★☆)"
```

## Threshold References

This library uses industry-standard thresholds from:

- **Powerful Signal** - Cellular signal booster manufacturer
- **Digi International** - Industrial cellular router manufacturer
- **Telco Antennas** - Professional antenna installation
- **3GPP TS 36.133** - Measurement ranges (operator-specific)
- **FreeRTOS Cellular Interface** - Implementation reference

### Default Thresholds

| Metric         | Excellent | Good        | Fair         | Poor   |
| -------------- | --------- | ----------- | ------------ | ------ |
| **RSRP** (dBm) | ≥ -89     | -90 to -104 | -105 to -114 | ≤ -115 |
| **RSRQ** (dB)  | ≥ -9      | -10 to -14  | -15 to -19   | ≤ -20  |
| **RSSI** (dBm) | ≥ -65     | -65 to -75  | -75 to -85   | ≤ -85  |
| **SINR** (dB)  | ≥ 13      | 6 to 13     | 0 to 6       | < 0    |

## Development

### Prerequisites

- [mise](https://mise.jdx.dev/) (task runner and tool manager)

### Running Tests

```bash
# Run all tests
mise run test

# Run CI checks (lint + test + coverage)
mise run ci

# Check coverage
mise run covercheck
```

## Documentation

See [pkg.go.dev](https://pkg.go.dev/github.com/hugoh/cellular-signal) for full API documentation.

## License

MIT License - see [LICENSE](LICENSE) for details.
