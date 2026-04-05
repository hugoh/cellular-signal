package signal

// defaultRSRPThresholds returns industry-standard RSRP thresholds.
// Sources and references:
//   - Powerful Signal (cellular signal booster manufacturer)
//   - Digi International (industrial cellular router manufacturer)
//   - Telco Antennas (professional antenna installation)
//   - 3GPP TS 36.133 defines measurement ranges (operator-specific thresholds)
//   - FreeRTOS Cellular Interface implementation
//
// Typical ranges: -44 dBm (excellent) to -140 dBm (no signal).
func defaultRSRPThresholds() []Threshold {
	return []Threshold{
		{MinValue: -89, MaxValue: 0, Quality: QualityExcellent},
		{MinValue: -104, MaxValue: -89, Quality: QualityGood},
		{MinValue: -114, MaxValue: -104, Quality: QualityFair},
		{MinValue: -124, MaxValue: -114, Quality: QualityPoor},
		{MinValue: -200, MaxValue: -124, Quality: QualityNone},
	}
}

// defaultRSRQThresholds returns industry-standard RSRQ thresholds.
// Sources and references:
//   - Powerful Signal
//   - 3GPP TS 36.133 defines measurement ranges (-43 dB to 20 dB for 5G)
//   - Industry practice: higher (less negative) values indicate better quality
//
// RSRQ = N * (RSRP / RSSI) where N is the number of Resource Blocks.
// Typical ranges: -3 dB (excellent) to -20 dB (poor).
func defaultRSRQThresholds() []Threshold {
	return []Threshold{
		{MinValue: -9, MaxValue: 20, Quality: QualityExcellent},
		{MinValue: -14, MaxValue: -9, Quality: QualityGood},
		{MinValue: -19, MaxValue: -14, Quality: QualityFair},
		{MinValue: -50, MaxValue: -19, Quality: QualityPoor},
	}
}

// defaultRSSIThresholds returns industry-standard RSSI thresholds.
// Sources and references:
//   - Digi International
//   - ESP-IDF WiFi RSSI thresholds
//   - RSSI is less commonly used in LTE/5G compared to RSRP
//
// Typical ranges: -50 dBm (excellent) to -110 dBm (poor).
func defaultRSSIThresholds() []Threshold {
	return []Threshold{
		{MinValue: -65, MaxValue: 0, Quality: QualityExcellent},
		{MinValue: -75, MaxValue: -65, Quality: QualityGood},
		{MinValue: -85, MaxValue: -75, Quality: QualityFair},
		{MinValue: -120, MaxValue: -85, Quality: QualityPoor},
	}
}

// defaultSINRThresholds returns industry-standard SINR thresholds.
// Sources and references:
//   - Powerful Signal
//   - Nature Scientific Reports (2026)
//   - Higher values indicate cleaner signal with less noise
//
// SINR is a positive number in most contexts; negative values indicate
// the signal is weaker than the noise floor.
func defaultSINRThresholds() []Threshold {
	return []Threshold{
		{MinValue: 13, MaxValue: 100, Quality: QualityExcellent},
		{MinValue: 6, MaxValue: 13, Quality: QualityGood},
		{MinValue: 0, MaxValue: 6, Quality: QualityFair},
		{MinValue: -100, MaxValue: 0, Quality: QualityPoor},
	}
}
