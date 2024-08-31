package util

const (
	USD = "USD"
	INR = "INR"
	CAD = "CAD"
)

func ValidateCurrency(s string) bool {
	switch s {
	case USD, INR, CAD:
		return true
	}
	return false
}
