# cellular-signal

[![Go Reference](https://pkg.go.dev/badge/github.com/hugoh/cellular-signal/v2.svg)](https://pkg.go.dev/github.com/hugoh/cellular-signal/v2)
[![CI](https://github.com/hugoh/cellular-signal/actions/workflows/ci.yml/badge.svg)](https://github.com/hugoh/cellular-signal/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/hugoh/cellular-signal/graph/badge.svg?token=UMZMODZ5PV)](https://codecov.io/github/hugoh/cellular-signal)
[![Go Report Card](https://goreportcard.com/badge/github.com/hugoh/cellular-signal/v2)](https://goreportcard.com/report/github.com/hugoh/cellular-signal/v2)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhugoh%2Fcellular-signal.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhugoh%2Fcellular-signal?ref=badge_shield)

Go library for rating cellular signal quality (LTE/4G/5G) based on industry standards.

## Features

- Rate RSRP, RSRQ, RSSI, and SINR signal metrics
- Based on industry-standard thresholds from telecom vendors
- Zero dependencies (pure Go stdlib)
- Fully tested
- Customizable thresholds

## Installation

```bash
go get github.com/hugoh/cellular-signal/v2
```

## Quick Start

```go
package main

import (
    "fmt"

    signal "github.com/hugoh/cellular-signal/v2"
)

func main() {
    rater := signal.NewRater()

    // Rate individual metrics
    rsrpRating := rater.RateRSRP(-92)
    fmt.Println(rsrpRating)
    // Output: RSRP: -92 dBm (Good ★★★★☆)

    // Access rating details
    fmt.Printf("Quality: %s\n", rsrpRating.Quality.String())
    fmt.Printf("Metric: %s\n", rsrpRating.Metric)
    fmt.Printf("Value: %v %s\n", rsrpRating.Value, rsrpRating.Metric.Unit())
}
```

## API Reference

### Creating a Rater

```go
// Default rater with industry-standard thresholds
rater := signal.NewRater()

// Custom thresholds: each entry is the lower bound of a quality level,
// ordered from best quality to worst (strictly descending MinValue).
customThresholds := []signal.Threshold{
    {MinValue: -80, Quality: signal.QualityExcellent},
    {MinValue: -100, Quality: signal.QualityGood},
    {MinValue: -200, Quality: signal.QualityPoor},
}
rater, err := signal.NewRaterWithThresholds(
    signal.WithRSRPThresholds(customThresholds),
)
if err != nil {
    log.Fatalf("Failed to create rater: %v", err)
}
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
formatted := rating.String() // Rating implements fmt.Stringer
// Output: "RSRP: -92 dBm (Good ★★★★☆)"

// Custom format
custom := rating.Format("%m=%v%u %s")
// Output: "RSRP=-92dBm ★★★★☆"
```

#### Format Verbs

| Verb | Description | Example |
| ---- | ----------- | ------- |
| `%m` | Metric name | `RSRP`  |
| `%v` | Value       | `-92`   |
| `%u` | Unit        | `dBm`   |
| `%q` | Quality     | `Good`  |
| `%s` | Stars       | `★★★★☆` |
| `%%` | Literal `%` | `%`     |

## Threshold References

This library uses industry-standard thresholds from:

- **Powerful Signal** - Cellular signal booster manufacturer
- **Digi International** - Industrial cellular router manufacturer
- **Telco Antennas** - Professional antenna installation
- **3GPP TS 36.133** - Measurement ranges (operator-specific)
- **FreeRTOS Cellular Interface** - Implementation reference

### Default Thresholds

Each level applies from its bound up to (but not including) the next
better level's bound.

| Metric         | Excellent | Good   | Fair   | Poor   | No Signal |
| -------------- | --------- | ------ | ------ | ------ | --------- |
| **RSRP** (dBm) | ≥ -89     | ≥ -104 | ≥ -114 | ≥ -124 | < -124    |
| **RSRQ** (dB)  | ≥ -9      | ≥ -14  | ≥ -19  | < -19  | —         |
| **RSSI** (dBm) | ≥ -65     | ≥ -75  | ≥ -85  | < -85  | —         |
| **SINR** (dB)  | ≥ 13      | ≥ 6    | ≥ 0    | < 0    | —         |

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

See [pkg.go.dev](https://pkg.go.dev/github.com/hugoh/cellular-signal/v2) for full API documentation.

## License

MIT License - see [LICENSE](LICENSE) for details.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhugoh%2Fcellular-signal.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhugoh%2Fcellular-signal?ref=badge_large)
