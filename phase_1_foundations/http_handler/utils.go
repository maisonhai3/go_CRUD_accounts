package httphandler

var validCurrencies = map[string]bool{
	"VND": true,
	"USD": true,
}

func isValidCurrency(currency string) bool {
	return validCurrencies[currency]
}
