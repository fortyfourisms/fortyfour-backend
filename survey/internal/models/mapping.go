package models

// IMPACT MAPPING
// int → string (ke database)
func MapImpactIntToString(i ImpactLevel) string {
	switch i {
	case ImpactNotSignificant:
		return "Tidak Signifikan"
	case ImpactFairlySignificant:
		return "Cukup Signifikan"
	case ImpactSignificant:
		return "Signifikan"
	case ImpactVerySignificant:
		return "Sangat Signifikan"
	default:
		return ""
	}
}

// string → int (dari database)
func MapImpactStringToInt(s string) ImpactLevel {
	switch s {
	case "Tidak Signifikan":
		return ImpactNotSignificant
	case "Sukup Signifikan":
		return ImpactFairlySignificant
	case "Signifikan":
		return ImpactSignificant
	case "Sangat Signifikan":
		return ImpactVerySignificant
	default:
		return 0
	}
}

// FREQUENCY MAPPING
// int → string
func MapFrequencyIntToString(f FrequencyLevel) string {
	switch f {
	case FrequencySmall:
		return "Kecil"
	case FrequencyMedium:
		return "Sedang"
	case FrequencyLarge:
		return "Besar"
	case FrequencyVeryLarge:
		return "Sangat Besar"
	default:
		return ""
	}
}

// string → int
func MapFrequencyStringToInt(s string) FrequencyLevel {
	switch s {
	case "Kecil":
		return FrequencySmall
	case "Sedang":
		return FrequencyMedium
	case "Besar":
		return FrequencyLarge
	case "Sangat Besar":
		return FrequencyVeryLarge
	default:
		return 0
	}
}
