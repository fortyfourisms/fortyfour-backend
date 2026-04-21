package models

// IMPACT MAPPING
// int → string (ke database)
func MapImpactIntToString(i ImpactLevel) string {
	switch i {
	case ImpactNotSignificant:
		return "tidak_signifikan"
	case ImpactFairlySignificant:
		return "cukup_signifikan"
	case ImpactSignificant:
		return "signifikan"
	case ImpactVerySignificant:
		return "sangat_signifikan"
	default:
		return ""
	}
}

// string → int (dari database)
func MapImpactStringToInt(s string) ImpactLevel {
	switch s {
	case "tidak_signifikan":
		return ImpactNotSignificant
	case "cukup_signifikan":
		return ImpactFairlySignificant
	case "signifikan":
		return ImpactSignificant
	case "sangat_signifikan":
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
		return "kecil"
	case FrequencyMedium:
		return "sedang"
	case FrequencyLarge:
		return "besar"
	case FrequencyVeryLarge:
		return "sangat_besar"
	default:
		return ""
	}
}

// string → int
func MapFrequencyStringToInt(s string) FrequencyLevel {
	switch s {
	case "kecil":
		return FrequencySmall
	case "sedang":
		return FrequencyMedium
	case "besar":
		return FrequencyLarge
	case "sangat_besar":
		return FrequencyVeryLarge
	default:
		return 0
	}
} 