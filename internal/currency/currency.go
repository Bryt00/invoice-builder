package currency

import "strings"

type CurrencyOption struct {
	Code   string
	Symbol string
	Name   string
	Flag   string
}

var supportedMap = map[string]CurrencyOption{
	"USD": {Code: "USD", Symbol: "$", Name: "US Dollar", Flag: "🇺🇸"},
	"EUR": {Code: "EUR", Symbol: "€", Name: "Euro", Flag: "🇪🇺"},
	"GBP": {Code: "GBP", Symbol: "£", Name: "British Pound", Flag: "🇬🇧"},
	"GHS": {Code: "GHS", Symbol: "GH₵", Name: "Ghanaian Cedi", Flag: "🇬🇭"},
	"NGN": {Code: "NGN", Symbol: "₦", Name: "Nigerian Naira", Flag: "🇳🇬"},
	"CAD": {Code: "CAD", Symbol: "CA$", Name: "Canadian Dollar", Flag: "🇨🇦"},
	"AUD": {Code: "AUD", Symbol: "A$", Name: "Australian Dollar", Flag: "🇦🇺"},
	"JPY": {Code: "JPY", Symbol: "¥", Name: "Japanese Yen", Flag: "🇯🇵"},
	"INR": {Code: "INR", Symbol: "₹", Name: "Indian Rupee", Flag: "🇮🇳"},
	"ZAR": {Code: "ZAR", Symbol: "R", Name: "South African Rand", Flag: "🇿🇦"},
	"BRL": {Code: "BRL", Symbol: "R$", Name: "Brazilian Real", Flag: "🇧🇷"},
	"AED": {Code: "AED", Symbol: "AED", Name: "UAE Dirham", Flag: "🇦🇪"},
}

func Symbol(code string) string {
	c := strings.ToUpper(strings.TrimSpace(code))
	if opt, ok := supportedMap[c]; ok {
		return opt.Symbol
	}
	if c == "" {
		return "$"
	}
	return c
}

func SupportedCurrencies() []CurrencyOption {
	return []CurrencyOption{
		supportedMap["USD"],
		supportedMap["EUR"],
		supportedMap["GBP"],
		supportedMap["GHS"],
		supportedMap["NGN"],
		supportedMap["CAD"],
		supportedMap["AUD"],
		supportedMap["JPY"],
		supportedMap["INR"],
		supportedMap["ZAR"],
		supportedMap["BRL"],
		supportedMap["AED"],
	}
}
