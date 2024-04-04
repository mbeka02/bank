package utils

const (
	KSH = "KSH"
	USD = "USD"
	EUR = "EUR"
)

func validateCurrency(currency string) bool {
	switch currency {
	case USD, EUR, KSH:
		return true
	}
	return false
}
