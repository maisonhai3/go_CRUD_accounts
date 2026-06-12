package handlers

func isValidCurrency(cur string) bool {
	// Whitelist
	switch cur {
	case "VND", "USD":
		return true
	}
	return false
}
